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

func AddChangePasswordCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "change-password --username [username] --new-password [newPassword]",
		Short:   "Change a password of a user",
		Long:    `This command changes the password of a user in your database instance.`,
		Example: `dbac db change-password --username engin --new-password newStrongPass123`,
		Args:    cobra.NoArgs,
		Run:     runChangePassword,
	}

	cmd.Flags().String("username", "", "Username of the user whose password is to be changed")
	cmd.Flags().String("new-password", "", "New password for the user")
	subcommand.AddCommand(cmd)
}

func runChangePassword(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	username, _ := cmd.Flags().GetString("username")
	newPassword, _ := cmd.Flags().GetString("new-password")
	if username == "" || newPassword == "" {
		fmt.Fprintf(os.Stderr, "ERROR: both --username and --new-password flags are required\n")
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
		psql.ChangeUserPassword(username, newPassword)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.ChangeUserPassword(username, newPassword)
		mysql.Close()
	}
}
