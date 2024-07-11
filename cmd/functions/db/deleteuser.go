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

func AddDeleteUserCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "delete-user --username [username]",
		Short:   "Delete a user",
		Long:    `This command deletes a user from your database instance.`,
		Example: "dbac db delete-user --username rumeysa",
		Args:    cobra.NoArgs,
		Run:     runDeleteUser,
	}
	cmd.Flags().String("username", "", "Username of the user to be deleted")
	subcommand.AddCommand(cmd)
}

func runDeleteUser(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	username, _ := cmd.Flags().GetString("username")
	if username == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --username flag is required\n")
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
		psql.DeleteUser(username)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.DeleteUser(username)
		mysql.Close()
	}
}
