package cmd

import (
	functions "dbac/cmd/functions/db"

	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage your databases",
	Long:  `Manage database operations such as user and grant operations.`,
	Example: `
# List all databases
dbac db list-databases

# Create a new database
dbac db create-database --name exampleDB

# Delete a database
dbac db delete-database --name exampleDB

# Ping the database to check connectivity
dbac db ping
`,
}

func init() {
	rootCmd.AddCommand(databaseCmd)

	functions.AddPingCommand(databaseCmd)
	functions.AddListDatabasesCommand(databaseCmd)
	functions.AddCreateDatabaseCommand(databaseCmd)
	functions.AddDeleteDatabaseCommand(databaseCmd)
	functions.AddListTablesCommand(databaseCmd)
	functions.AddListUsersCommand(databaseCmd)
	functions.AddCreateUserCommand(databaseCmd)
	functions.AddDeleteUserCommand(databaseCmd)
	functions.AddChangePasswordCommand(databaseCmd)
	functions.AddGrantDatabaseCommand(databaseCmd)
	functions.AddRevokeDatabaseCommand(databaseCmd)
	functions.AddGrantTableCommand(databaseCmd)
	functions.AddRevokeTableCommand(databaseCmd)
	functions.AddExecCommand(databaseCmd)
	// functions.AddDumpCommand(databaseCmd)
}
