package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fermyon/verman-plugin/internal/verman"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets Spin to the requested version.",
	Long:  "Sets Spin to the requested version. If the requested version is not found locally and exists in the remote repository, it will be downloaded.",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := verman.GetDesiredVersionForSet(args)
		if err != nil {
			return err
		}

		versionDir, err := getVersionDir()
		if err != nil {
			return err
		}

		symlinkDir := path.Join(versionDir, "current_version")
		binaryDir := path.Join(versionDir, version)

		if err := checkPathVar(symlinkDir); err != nil {
			return err
		}

		if err := downloadSpin(versionDir, version); err != nil {
			return err
		}

		if err = updateSpinBinary(binaryDir, symlinkDir); err != nil {
			return err
		}

		fmt.Printf("Spin has been updated to version %s\n", version)
		return nil
	},
}

var setLatestStableCmd = &cobra.Command{
	Use:   "latest",
	Short: "Sets Spin to the latest stable version",
	Long:  "Sets Spin to the latest stable version and will download the stable version binary if not found locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		versionDir, err := getVersionDir()
		if err != nil {
			return err
		}

		version, err := getLatestTag()
		if err != nil {
			return err
		}

		symlinkDir := path.Join(versionDir, "current_version")
		binaryDir := path.Join(versionDir, version)

		if err := checkPathVar(symlinkDir); err != nil {
			return err
		}

		if err := downloadSpin(versionDir, version); err != nil {
			return err
		}

		if err := updateSpinBinary(binaryDir, symlinkDir); err != nil {
			return err
		}

		fmt.Println("Spin has been updated to the latest stable version")
		return nil
	},
}

// updateSpinBinary creates a symlink pointing to a binary file containing the specified version of Spin
func updateSpinBinary(binaryDir, symlinkDir string) error {
	if err := os.MkdirAll(symlinkDir, 0755); err != nil {
		return err
	}

	// If there is already an existing symlink, this deletes the symlink (or deletes nothing if the symlink doesn't exist) so a new one can be created
	if err := os.Remove(path.Join(symlinkDir, "spin")); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove old symlink: %v", err)
		}
	}

	if err := os.Symlink(path.Join(binaryDir, "spin"), path.Join(symlinkDir, "spin")); err != nil {
		return err
	}

	return nil
}

func checkPathVar(dirPath string) error {
	// Check to make sure the currentVersionPath is in the $PATH variable
	path := os.Getenv("PATH")
	pathSeparator := string(os.PathListSeparator)
	pathIsInPATH := false
	for _, p := range strings.Split(path, pathSeparator) {
		if p == dirPath {
			pathIsInPATH = true
		}
	}
	if !pathIsInPATH {
		return fmt.Errorf("unable to find %q in $PATH", dirPath)
	}

	return nil
}
