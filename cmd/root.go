package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "verman",
	Short: "A plugin for Spin that makes it easy to manage different versions of the Spin CLI.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Set
	setCmd.AddCommand(setLatestStableCmd)
	rootCmd.AddCommand(setCmd)
	//Alias
	rootCmd.AddCommand(aliasCmd)
	// Get
	getCmd.AddCommand(getLatestStableCmd)
	rootCmd.AddCommand(getCmd)
	// List
	rootCmd.AddCommand(listCmd)
	// List Remote
	rootCmd.AddCommand(listRemoteCmd)
	// Remove
	removeCmd.AddCommand(removeAllCmd)
	removeCmd.AddCommand(removeCurrentCmd)
	rootCmd.AddCommand(removeCmd)
	// Update
	updateCmd.AddCommand(updateCanaryCmd)
	rootCmd.AddCommand(updateCmd)
}
