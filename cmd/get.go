package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/fermyon/verman-plugin/internal/verman"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Downloads the binary for the requested version if not found locally.",
	Long:  "Downloads the binary for the requested version if not found locally. Multiple versions can be downloaded at once: \"spin verman get 2.1.0 canary\".",
	RunE: func(cmd *cobra.Command, args []string) error {
		versions, err := verman.GetDesiredVersionsForGet(args)
		if err != nil {
			return err
		}

		versionDir, err := getVersionDir()
		if err != nil {
			return err
		}

		for _, version := range versions {
			if !strings.HasPrefix(version, "v") && version != "canary" {
				version = "v" + version
			}

			if err := downloadSpin(versionDir, version); err != nil {
				return err
			}
		}

		return nil
	},
}

var getLatestStableCmd = &cobra.Command{
	Use:   "latest",
	Short: "Downloads the binary for the latest stable version if not found locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		versionDir, err := getVersionDir()
		if err != nil {
			return err
		}

		version, err := getLatestTag()
		if err != nil {
			return err
		}

		if err := downloadSpin(versionDir, version); err != nil {
			return err
		}

		return nil
	},
}

// exists indicates whether the file/directory path exists
func exists(path string) (bool, error) {
	// If the path does exist...
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	// If the path doesn't exist...
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// getVersionDir returns the directory in which the "spin verman" version files will be stored
func getVersionDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	versionDir := path.Join(homeDir, ".spin_verman", "versions")

	dirExists, err := exists(versionDir)
	if err != nil {
		return "", err
	}

	if !dirExists {
		if err = os.MkdirAll(versionDir, 0755); err != nil {
			return "", err
		}

		return "", nil
	}

	return versionDir, nil
}

// getLatestTag returns a string containing the tag of the latest stable version of Spin
func getLatestTag() (string, error) {
	var latestRelease struct {
		TagName string `json:"tag_name"`
	}

	resp, err := http.Get("https://api.github.com/repos/fermyon/spin/releases/latest")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve the tag for the latest stable version of Spin: %v", err)
	}

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(jsonBytes, &latestRelease); err != nil {
		return "", err
	}

	return latestRelease.TagName, nil
}

// downloadSpin will retrieve the desired version of Spin if it is not present in the version directory
func downloadSpin(versionDir, version string) error {
	var spinArch string
	var spinOS string

	// Checking for compatible architectures
	if runtime.GOARCH == "amd64" {
		spinArch = "amd64"
	} else if runtime.GOARCH == "arm64" {
		spinArch = "aarch64"
	} else {
		return fmt.Errorf("%q is not an architecture that Spin supports", runtime.GOARCH)
	}

	// Checking for compatible operating systems
	if runtime.GOOS == "linux" {
		// TODO: When would we want to download 'static-linux' vs just 'linux'?
		spinOS = "linux"
	} else if runtime.GOOS == "darwin" {
		spinOS = "macos"
	} else {
		return fmt.Errorf("%q is not an OS that this Spin plugin supports", runtime.GOOS)
	}

	if !isSemver(version) && version != "canary" {
		return fmt.Errorf("the requested version %q is not proper semver (i.e. v0.0.0 or 0.0.0)", version)
	}

	fileName := fmt.Sprintf("spin-%s-%s-%s.tar.gz", version, spinOS, spinArch)

	dirExists, err := exists(versionDir)
	if err != nil {
		return err
	}

	// Determines if we need to pull the file from GitHub
	var versionFolderExists bool

	if !dirExists {
		if err = os.MkdirAll(versionDir, 0755); err != nil {
			return err
		}
	} else {
		dirFiles, err := os.ReadDir(versionDir)
		if err != nil {
			return err
		}

		for _, file := range dirFiles {
			// Checking if the Spin binary has previously been unpacked...
			if file.Name() == version {
				fmt.Printf("Spin version %s found locally.\n", version)
				versionFolderExists = true
				break
			}
		}
	}

	// If the tar.gz file doesn't exist, pull from GitHub
	if !versionFolderExists {
		fmt.Printf("Spin version %s not found locally. Retrieving from source...\n", version)

		resp, err := http.Get("https://github.com/fermyon/spin/releases/download/" + version + "/" + fileName)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("the version number provided is invalid: %s", version)
		}

		out, err := os.Create(path.Join(versionDir, fileName))
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("Spin version %s was retrieved successfully!\n", version)
		if err = unpackSpin(versionDir, fileName, version); err != nil {
			return err
		}
	}

	return nil
}

// unpackSpin unpacks the binary file from a .tar.gz file for the specified version of Spin
func unpackSpin(directory, tarGzFileName, version string) error {
	if err := os.Chdir(directory); err != nil {
		return err
	}

	gzipStream, err := os.ReadFile(tarGzFileName)
	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(bytes.NewReader(gzipStream))
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("unpackSpin: Next() failed: %w", err)
		}

		// Extracting only the Spin CLI binary
		if header.Typeflag == tar.TypeReg && header.Name == "spin" {
			// Create the file with the original permissions
			outFile, err := os.OpenFile(header.Name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			outFile.Close()

			// Ensure the file has the correct permissions
			if err := os.Chmod(header.Name, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("unpackSpin: could not set file permissions: %w", err)
			}
		}
	}

	// Create a folder named with the relevant Spin version
	if err := os.MkdirAll(version, 0755); err != nil {
		return err
	}

	if err := os.Rename("spin", path.Join(directory, version, "spin")); err != nil {
		return err
	}

	if err := os.Remove(tarGzFileName); err != nil {
		return err
	}

	return nil
}

// isSemver makes sure the version passed is proper semver
func isSemver(version string) bool {
	return semver.IsValid(version)
}
