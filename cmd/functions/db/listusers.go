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

func AddListUsersCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "list-users",
		Short:   "List users",
		Long:    `This command lists users in the currently selected database instance.`,
		Example: "dbac db list-users",
		Args:    cobra.NoArgs,
		Run:     runListUsers,
	}
	subcommand.AddCommand(cmd)
}

func runListUsers(cmd *cobra.Command, args []string) {
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
		psql.ListUsers()
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.ListUsers()
		mysql.Close()
	}
}
