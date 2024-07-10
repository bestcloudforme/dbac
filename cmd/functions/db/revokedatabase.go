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

func AddRevokeDatabaseCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "revoke-database",
		Short: "Revoke database permissions from a user",
		Long:  `This command revokes specified database permissions from a user.`,
		Example: `
# Revoke SELECT permission from a user on a specific database
dbac db revoke-database --username gorkem --permission SELECT --database mydb

# Revoke ALL permissions from a user on a specific database
dbac db revoke-database --username gorkem --permission ALL --database mydb
`,
		Args: cobra.NoArgs,
		Run:  runRevokeDatabase,
	}
	cmd.Flags().String("username", "", "Username of the user")
	cmd.Flags().String("permission", "", "Permissions to revoke")
	cmd.Flags().String("database", "", "Database from which permissions are revoked")
	subcommand.AddCommand(cmd)
}

func runRevokeDatabase(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	username, _ := cmd.Flags().GetString("username")
	permission, _ := cmd.Flags().GetString("permission")
	database, _ := cmd.Flags().GetString("database")
	if username == "" || permission == "" || database == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --username, --permission, and --database flags are required\n")
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
		psql.RevokePermissions(database, username, permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.RevokePermissions(database, username, permission)
		mysql.Close()
	}
}
