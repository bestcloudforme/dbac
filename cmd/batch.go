package cmd

import (
	"dbac/cmd/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var batchCmd = &cobra.Command{
	Use:     "batch",
	Short:   "Run DBAC Commands with Batch Operations",
	Long:    `Runs batch operations for DBAC commands from a specified file.`,
	Example: `dbac batch --file "/path/to/batchfile.yml"`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR retrieving file flag: %v\n", err)
			os.Exit(1)
		}
		if file == "" {
			fmt.Fprintf(os.Stderr, "ERROR: --file flag is required\n")
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to display help: %v\n", err)
			}
			os.Exit(1)
		}
		helper.RunBatch(file)
	},
}

func init() {
	batchCmd.Flags().StringP("file", "f", "", "File containing batch operations to execute")
	rootCmd.AddCommand(batchCmd)
}
