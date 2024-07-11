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

func AddGrantDatabaseCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "grant-database",
		Short: "Grant permissions on a database to a user",
		Long:  `This command grants specified permissions to a user for a database.`,
		Example: `
# Grant SELECT permission on a database to a user
dbac db grant-database --username gorkem --permission SELECT --database mydb

# Grant ALL permissions on a database to a user
dbac db grant-database --username gorkem --permission ALL --database mydb
`,
		Args: cobra.NoArgs,
		Run:  runGrantDatabase,
	}
	cmd.Flags().String("username", "", "Username of the user")
	cmd.Flags().String("permission", "", "Permissions to grant")
	cmd.Flags().String("database", "", "Database on which permissions are granted")
	subcommand.AddCommand(cmd)
}

func runGrantDatabase(cmd *cobra.Command, args []string) {
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
		psql.GrantPermissions(database, username, permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
		mysql.GrantPermissions(database, username, permission)
		mysql.Close()
	}
}
