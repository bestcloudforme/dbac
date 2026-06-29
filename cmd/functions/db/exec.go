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
		case "psql":
			dbPort, _ := strconv.Atoi(profile.Port)
			if err := psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database, profile.SSLMode); err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
				os.Exit(1)
			}
			if err := psql.FileExec(file); err != nil {
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
			if err := mysql.FileExec(file); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if err := mysql.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing connection: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		switch profile.DbType {
		case "psql":
			dbPort, _ := strconv.Atoi(profile.Port)
			if err := psql.NewConnection(profile.Host, dbPort, profile.User, profile.Password, profile.Database, profile.SSLMode); err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
				os.Exit(1)
			}
			if err := psql.Exec(query); err != nil {
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
			if err := mysql.Exec(query); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if err := mysql.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing connection: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
