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

func AddListDatabasesCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "list-databases",
		Short:   "List your databases",
		Long:    `This command lists all databases in the currently selected database instance.`,
		Example: "dbac db list-databases",
		Args:    cobra.NoArgs,
		Run:     runListDatabases,
	}
	subcommand.AddCommand(cmd)
}

func runListDatabases(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	profile := helper.ReadProfile(currentProfileName)
	switch profile.DbType {
	case "postgres":
		dbPort, _ := strconv.Atoi(profile.Port)
		psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database)
		psql.ListDatabases()
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.ListDatabases()
		mysql.Close()
	}
}
