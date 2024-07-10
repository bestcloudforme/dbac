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

func AddExecCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute SQL commands on your database",
		Long:  `This command executes SQL commands or scripts on your database instance.`,
		Example: `
# Execute a single SQL query directly:
dbac db exec --query "SELECT * FROM users"

# Execute SQL commands from a file:
dbac db exec --file "/path/to/your/sqlfile.sql"
`,
		Args: cobra.NoArgs,
		Run:  runExecCommand,
	}
	cmd.Flags().StringP("file", "f", "", "File containing SQL commands to execute")
	cmd.Flags().String("query", "", "SQL query to execute")
	subcommand.AddCommand(cmd)
}

func runExecCommand(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}
	file, _ := cmd.Flags().GetString("file")
	query, _ := cmd.Flags().GetString("query")
	if query == "" && file == "" {
		fmt.Fprintf(os.Stderr, "ERROR: either --query or --file flag is required\n")
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to display help: %v\n", err)
		}
		os.Exit(1)
	}
	profile := helper.ReadProfile(currentProfileName)
	if file != "" {
		switch profile.DbType {
		case "postgres":
			dbPort, _ := strconv.Atoi(profile.Port)
			psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database)
			psql.FileExec(file)
			psql.Close()
		case "mysql":
			mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
			mysql.FileExec(file)
			mysql.Close()
		}
	} else {
		switch profile.DbType {
		case "postgres":
			dbPort, _ := strconv.Atoi(profile.Port)
			psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database)
			psql.Exec(query)
			psql.Close()
		case "mysql":
			mysql.NewConnection(profile.Host, profile.Port, profile.User, profile.Password, profile.Database)
			mysql.Exec(query)
			mysql.Close()
		}
	}
}
