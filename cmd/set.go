package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets Spin to the requested version.",
	Long:  "Sets Spin to the requested version, and will download the binary for the requested version if not found locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must indicate the version of Spin you wish to set")
		}

		version := args[0]

		if !strings.HasPrefix(version, "v") && version != "canary" {
			version = "v" + version
		}

		versionDir, err := getVersionDir()
		if err != nil {
			return err
		}

		if err := downloadSpin(versionDir, version); err != nil {
			return err
		}

		if err = updateSpinBinary(versionDir, version); err != nil {
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

		if err := downloadSpin(versionDir, version); err != nil {
			return err
		}

		if err := updateSpinBinary(versionDir, version); err != nil {
			return err
		}

		fmt.Println("Spin has been updated to the latest stable version")
		return nil
	},
}

// updateSpinBinary creates a symlink pointing to a binary file containing the specified version of Spin
func updateSpinBinary(directory, version string) error {
	if err := os.MkdirAll(path.Join(directory, "current_version"), 0755); err != nil {
		return err
	}

	symLinkDir := path.Join(directory, "current_version")

	// Removing old SymLink, returning an error only if the error is not a 'file does not exist' error
	if err := os.Remove(path.Join(symLinkDir, "spin")); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove old symlink: %v", err)
		}
	}

	if err := os.Symlink(path.Join(directory, version, "spin"), path.Join(symLinkDir, "spin")); err != nil {
		return err
	}

	testSpinVersionCmd := exec.Command("spin", "--version")
	currentSpinVersionBytes, err := testSpinVersionCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error getting the current version of Spin: %v\n%s", err, string(currentSpinVersionBytes))
	}

	if version == "canary" {
		// Checking the version of the canary
		canaryFile := path.Join(directory, "canary", "spin")
		cmd := exec.Command(canaryFile, "--version")
		canaryVersionBytes, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error getting the current canary version: %v\n%s", err, string(canaryVersionBytes))
		}

		// Retrieves the version number from the "spin --version" command
		version = strings.Split(string(canaryVersionBytes), " ")[1]

	} else {
		// Remove the "v" prefix from the version
		version = version[1:]
	}

	if !strings.Contains(string(currentSpinVersionBytes), version) {
		return fmt.Errorf("it looks like the version of the current Spin executable does not match what was requested, so please check to make sure the path %q is prepended to your path", symLinkDir)
	}

	return nil
}
