package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbac",
	Short: "DBAC - Database Access Control CLI",
	Long: `DBAC is a command-line interface that simplifies 
the management of databases, users, and permissions. It provides a suite of tools 
for database operations, including creating, modifying, and deleting databases and 
users, as well as managing access controls.

For example, you can easily manage databases and their access levels with commands like:
  dbac create-user --username gorkem
  dbac grant-database --username gorkem --permission read --database testdb`,

	Example: `  dbac list-users
  dbac revoke-table --username rumeysa --table orders --permission write`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "DBAC CLI encountered an error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
}
