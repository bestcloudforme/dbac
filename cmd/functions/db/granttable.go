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

func AddGrantTableCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "grant-table",
		Short: "Grant table-specific permissions to a user",
		Long:  `This command grants specified table-specific permissions to a user.`,
		Example: `
# Grant SELECT permission on a specific table to a user
dbac db grant-table --username gorkem --permission SELECT --table employees

# Grant ALL permissions on a specific table to a user
dbac db grant-table --username gorkem --permission ALL --table employees
`,
		Args: cobra.NoArgs,
		Run:  runGrantTable,
	}
	cmd.Flags().String("username", "", "Username of the user")
	cmd.Flags().String("permission", "", "Permissions to grant")
	cmd.Flags().String("table", "", "Table on which permissions are granted")
	subcommand.AddCommand(cmd)
}

func runGrantTable(cmd *cobra.Command, args []string) {
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
		psql.GrantTablePermissions(table, username, permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.GrantTablePermissions(profile.Database, table, username, permission)
		mysql.Close()
	}
}
