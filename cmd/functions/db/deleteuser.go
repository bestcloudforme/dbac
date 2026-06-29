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
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
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
	yes, _ := cmd.Flags().GetBool("yes")
	if !yes && !helper.Confirm(fmt.Sprintf("Delete user %q?", username)) {
		fmt.Fprintln(os.Stderr, "Aborted.")
		os.Exit(1)
	}
	profile := helper.ReadProfile(currentProfileName)
	switch profile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(profile.Port)
		if err := psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database, profile.SSLMode); err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		if err := psql.DeleteUser(username); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := psql.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing connection: %v\n", err)
			os.Exit(1)
		}
	case "mysql":
		if err := mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database); err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		if err := mysql.DeleteUser(username); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := mysql.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing connection: %v\n", err)
			os.Exit(1)
		}
	}
}
