package functions

import (
	"dbac/cmd/helper"
	"fmt"

	_ "github.com/lib/pq"
)

var currentProfile helper.Profile

func App(params []string, profile helper.Profile) {
	currentProfile = profile

	if len(params) < 1 {
		printAppHelp()
		return
	}

	switch params[0] {
	case "batch":
		Batch(params)

	case "init":
		Init(params)

	case "profile":
		Profile(params)

	case "database":
		Database(params)

	case "-h":
		printAppHelp()

	default:
		fmt.Println("Invalid command.")
		printAppHelp()
	}
}

func printAppHelp() {
	fmt.Println("Available commands and options:")
	fmt.Println("  dbac init - Initialize the CLI configuration.")
	fmt.Println("  dbac profile - Manage profiles with options:")
	fmt.Println("    current - Display the current profile")
	fmt.Println("    list - List all profiles")
	fmt.Println("    switch [PROFILE-NAME] - Switch to a specified profile")
	fmt.Println("    add --db-type [DB-TYPE] --host [HOST] --user [USER] --port [PORT] --password [PASSWORD] --database [DATABASE] --profile-name [NAME] - Add a new profile with specified details")
	fmt.Println("    add --file [FILE] - Add a new profile from a file")
	fmt.Println("    delete --profile-name [PROFILE-NAME] - Delete a specified profile")
	fmt.Println("  dbac database - Perform database operations with options:")
	fmt.Println("    ping - Ping the current database to check the connection")
	fmt.Println("    list-user - List all users")
	fmt.Println("    list-databases - List all databases")
	fmt.Println("    list-tables - List all tables in the current database")
	fmt.Println("    create-user --username [USERNAME] --user-password [PASSWORD] - Create a new database user")
	fmt.Println("    create-database --database [DATABASE] - Create a new database")
	fmt.Println("    delete-user --username [USERNAME] - Delete an existing database user")
	fmt.Println("    delete-database --database [DATABASE] - Delete an existing database")
	fmt.Println("    change-password --username [USERNAME] --new-password [NEW_PASSWORD] - Change the password of an existing user")
	fmt.Println("    grant-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE] - Grant database-level permissions to a user")
	fmt.Println("    grant-table --username [USERNAME] --permission [PERMISSION] --table [TABLE] - Grant table-level permissions to a user")
	fmt.Println("    revoke-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE] - Revoke database-level permissions from a user")
	fmt.Println("    revoke-table --username [USERNAME] --permission [PERMISSION] --table [TABLE] - Revoke table-level permissions from a user")
	fmt.Println("    exec --query [QUERY] or --file [FILE] - Execute a SQL query or SQL commands from a file")
	fmt.Println("  dbac batch --file=\"batch.yaml\" - Execute batch commands from a YAML file.")
}
