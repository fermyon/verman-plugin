package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lists available Spin versions and aliases.",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := list()
		if err != nil {
			return err
		}

		fmt.Println(output)

		return nil
	},
}

// list prints all subdirectories (excluding current_version) in ~/.spin_verman/versions
func list() (string, error) {
	versionDir, err := getVersionDir()
	if err != nil {
		return "", err
	}

	pathExists, err := exists(versionDir)
	if err != nil {
		return "", err
	}

	if !pathExists {
		return "", nil
	}

	files, err := os.ReadDir(versionDir)
	if err != nil {
		return "", err
	}

	var output []string

	for _, file := range files {
		if file.Name() != "current_version" {
			output = append(output, file.Name())
		}
	}

	if len(output) == 0 {
		fmt.Println("No versions of Spin were found in the \"~/.spin_verman/versions\" directory. Run \"spin verman get --help\" to get started")
	}

	return strings.Join(output, "\n"), nil
}
