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

func AddRevokeTableCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "revoke-table",
		Short: "Revoke table-specific permissions from a user",
		Long:  `This command revokes specified table-specific permissions from a user.`,
		Example: `
# Revoke SELECT permission from a user on a specific table
dbac db revoke-table --username gorkem --permission SELECT --table employees

# Revoke ALL permissions from a user on a specific table
dbac db revoke-table --username gorkem --permission ALL --table employees
`,
		Args: cobra.NoArgs,
		Run:  runRevokeTable,
	}
	cmd.Flags().String("username", "", "Username of the user")
	cmd.Flags().String("permission", "", "Permissions to revoke")
	cmd.Flags().String("table", "", "Table from which permissions are revoked")
	subcommand.AddCommand(cmd)
}

func runRevokeTable(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	username, _ := cmd.Flags().GetString("username")
	permission, _ := cmd.Flags().GetString("permission")
	table, _ := cmd.Flags().GetString("table")
	if username == "" || permission == "" || table == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --username, --permission, and --table flags are required\n")
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
		psql.RevokeTablePermissions(table, username, permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.RevokeTablePermissions(profile.Database, table, username, permission)
		mysql.Close()
	}
}
