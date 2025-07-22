package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Removes the specified Spin version (or symlink if it's an alias) from the local directory.",
	Long:    "Removes the specified Spin version (or symlink if it's an alias) from the local directory. Only removes the relevant Spin binary located in the \"~/.spin_verman/versions\" directory.",
	Args:    cobra.MaximumNArgs(1), // This intentionally only removes one at a time
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must indicate which version of Spin you wish to delete or use one of the available subcommands")
		}

		version := args[0]

		if err := remove(version); err != nil {
			return err
		}

		return nil
	},
}

var removeCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Removes the alternate Spin version, reverting back to the root version of Spin.",
	Long:  "Removes the alternate Spin version, reverting back to the root version of Spin, but preserving all other versions of Spin downloaded locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := remove("current_version"); err != nil {
			return err
		}

		return nil
	},
}

var removeAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Removes all Spin versions (or symlinks if there are aliases) from the local directory.",
	Long:  "Removes all Spin versions (or symlinks if there are aliases) from the local directory. Only removes the Spin binaries located in the \"~/.spin_verman/versions\" directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Are you sure you want to delete all Spin versions?\nType \"y\", \"yes\", or any other key to cancel: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		output := strings.ToLower(input.Text())

		if output == "y" || output == "yes" {
			if err := removeAll(); err != nil {
				return err
			}

			fmt.Println("All Spin versions successfully deleted")

		} else {
			fmt.Println("No Spin versions were deleted")
		}

		return nil
	},
}

// remove removes the associated version file and directory. If the path isn't found, this will not return an error.
func remove(version string) error {
	versionDir, err := getVersionDir()
	if err != nil {
		return err
	}

	// Ensures that versions passed without a `v` prefix are deleted
	versionPath := path.Join(versionDir, version)
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		versionPath = path.Join(versionDir, "v"+version)
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			fmt.Println("Warning: file does not exist; nothing to remove")
			return nil
		}
		version = "v" + version
	}

	filePath := path.Join(versionDir, version)

	if err := os.RemoveAll(filePath); err != nil {
		return err
	}

	return nil
}

// removeAll removes all subdirectories in ~.spin_verman/versions
func removeAll() error {
	versionString, err := list()
	if err != nil {
		return err
	}

	for _, version := range strings.Split(versionString, "\n") {
		if err := remove(version); err != nil {
			return err
		}
	}

	// The list method doesn't return the "current_version" directory, so we need to manually delete it
	if err := remove("current_version"); err != nil {
		return err
	}

	return nil
}
