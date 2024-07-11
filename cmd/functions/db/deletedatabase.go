package functions

import (
	"dbac/cmd/helper"
	"dbac/cmd/mysql"
	"dbac/cmd/psql"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func AddDeleteDatabaseCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "delete-database --database [database_name]",
		Short:   "Delete a database",
		Long:    `This command deletes a database from your database instance.`,
		Example: "dbac db delete-database --database exampleDB",
		Args:    cobra.NoArgs,
		Run:     runDeleteDatabase,
	}
	cmd.Flags().String("database", "", "Database name to be deleted")
	subcommand.AddCommand(cmd)
}

func runDeleteDatabase(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	database, _ := cmd.Flags().GetString("database")
	if database == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --database flag is required\n")
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to display help: %v\n", err)
		}
		os.Exit(1)
	}
	profile := helper.ReadProfile(currentProfileName)
	switch profile.DbType {
	case "postgres":
		dbPort, _ := strconv.Atoi(profile.Port)
		psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database)
		psql.DeleteDatabase(database)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.DeleteDatabase(database)
		mysql.Close()
	}
}
