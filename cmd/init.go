package cmd

import (
	"dbac/cmd/helper"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize the application settings",
	Long:    `Sets up initial configuration files and directories required for the application to function properly.`,
	Example: `dbac init"`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		helper.InitProfile()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
