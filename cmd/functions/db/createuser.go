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

func AddCreateUserCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "create-user --username [username] --password [password]",
		Short:   "Create a user in database",
		Long:    `This command creates a user in your database instance.`,
		Example: "dbac db create-user --username engin --password newpassword123",
		Args:    cobra.NoArgs,
		Run:     runCreateUser,
	}
	cmd.Flags().String("username", "", "Username for the user")
	cmd.Flags().String("password", "", "Password for the user")
	subcommand.AddCommand(cmd)
}

func runCreateUser(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	if username == "" || password == "" {
		fmt.Fprintf(os.Stderr, "ERROR: both --username and --password flags are required\n")
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
		psql.CreateUser(username, password)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.CreateUser(username, password)
		mysql.Close()
	}
}
