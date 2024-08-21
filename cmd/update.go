package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the binary files found locally for the requested versions of Spin. Currently only supports the \"canary\" subcommand.",
}

var updateCanaryCmd = &cobra.Command{
	Use:   "canary",
	Short: "Updates the binary files found locally for the canary version of Spin. If the canary version is not found locally, it will be retrieved from source.",
	RunE: func(cmd *cobra.Command, args []string) error {
		versionDir, err := getVersionDir()
		if err != nil {
			return err
		}

		_, canaryFileErr := os.Stat(path.Join(versionDir, "canary"))
		if err != nil {
			if !os.IsNotExist(canaryFileErr) {
				return canaryFileErr
			}
		}

		if err := remove("canary"); err != nil {
			return err
		}

		// If the canary file already existed locally...
		if !os.IsNotExist(canaryFileErr) {
			fmt.Println("Old canary version successfully deleted")
		}

		if err := downloadSpin(versionDir, "canary"); err != nil {
			return err
		}

		return nil
	},
}
