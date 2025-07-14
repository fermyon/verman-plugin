package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias [name] [path]",
	Short: "Creates an alias for a local Spin binary.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("invalid arguments")
		}

		versionDir, err := getVersionDir()
		if err != nil {
			return err
		}

		alias, filePath := args[0], args[1]

		aliasPath := path.Join(versionDir, alias)

		if err := os.MkdirAll(aliasPath, 0755); err != nil {
			return err
		}

		// In the case that there is already an existing alias + symlink,
		// this deletes the symlink (or deletes nothing if the symlink doesn't exist) so a new one can be created
		if err := os.Remove(path.Join(aliasPath, "spin")); err != nil {
			if !os.IsNotExist(err) { // We don't care if we are deleting nothing, so this ignores any `os.IsNotExist` errors
				return fmt.Errorf("failed to remove old symlink: %v", err)
			}
		}

		if err := os.Symlink(filePath, path.Join(aliasPath, "spin")); err != nil {
			return err
		}

		fmt.Printf("Created alias %q", alias)

		return nil
	},
}
