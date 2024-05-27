package functions

import (
	"dbac/cmd/mysql"
	"dbac/cmd/psql"
	"flag"
	"fmt"
	"strconv"
)

func Database(params []string) {
	if len(params) < 2 {
		printDatabaseHelp()
		return
	}

	switch params[1] {
	case "ping":
		pingDatabase()

	case "list-user":
		listUsers()

	case "create-user":
		createUser(params[2:])

	case "create-database":
		createDatabase(params[2:])

	case "delete-user":
		deleteUser(params[2:])

	case "delete-database":
		deleteDatabase(params[2:])

	case "change-password":
		changePassword(params[2:])

	case "grant-database":
		grantDatabase(params[2:])

	case "grant-table":
		grantTable(params[2:])

	case "revoke-database":
		revokeDatabase(params[2:])

	case "revoke-table":
		revokeTable(params[2:])

	case "list-databases":
		listDatabases()

	case "list-tables":
		listTables()

	case "exec":
		execQuery(params[2:])

	case "dump":
		dumpDatabase(params[2:])

	case "-h":
		printDatabaseHelp()

	default:
		fmt.Println("Invalid command.")
		printDatabaseHelp()
	}
}

func pingDatabase() {
	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.Ping()
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.Ping()
		mysql.Close()
	}
}

func listUsers() {
	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.ListUsers()
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.ListUsers()
		mysql.Close()
	}
}

func createUser(params []string) {
	cmd := flag.NewFlagSet("create-user", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the user")
	userPassword := cmd.String("user-password", "", "Password for the new database user")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *username == "" || *userPassword == "" {
		fmt.Println("Usage: dbac database create-user --username [USERNAME] --user-password [PASSWORD]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.CreateUser(*username, *userPassword)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.CreateUser(*username, *userPassword)
		mysql.Close()
	}
}

func createDatabase(params []string) {
	cmd := flag.NewFlagSet("create-database", flag.ExitOnError)
	databaseName := cmd.String("database", "", "Database name to be created")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments for creating database: %v\n", err)
		return
	}

	if *databaseName == "" {
		fmt.Println("Usage: dbac database create-database --database [DATABASE]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.CreateDatabase(*databaseName)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.CreateDatabase(*databaseName)
		mysql.Close()
	}
}

func deleteUser(params []string) {
	cmd := flag.NewFlagSet("delete-user", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the user to be deleted")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *username == "" {
		fmt.Println("Usage: dbac database delete-user --username [USERNAME]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.DeleteUser(*username)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.DeleteUser(*username)
		mysql.Close()
	}
}

func deleteDatabase(params []string) {
	cmd := flag.NewFlagSet("delete-database", flag.ExitOnError)
	databaseName := cmd.String("database", "", "Database name to be deleted")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *databaseName == "" {
		fmt.Println("Usage: dbac database delete-database --database [DATABASE]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.DeleteDatabase(*databaseName)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.DeleteDatabase(*databaseName)
		mysql.Close()
	}
}

func changePassword(params []string) {
	cmd := flag.NewFlagSet("change-password", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the user")
	newPassword := cmd.String("new-password", "", "New password of the user")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *username == "" || *newPassword == "" {
		fmt.Println("Usage: dbac database change-password --username [USERNAME] --new-password [NEW_PASSWORD]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.ChangeUserPassword(*username, *newPassword)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.ChangeUserPassword(*username, *newPassword)
		mysql.Close()
	}
}

func grantDatabase(params []string) {
	cmd := flag.NewFlagSet("grant-database", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the user")
	permission := cmd.String("permission", "", "Permission of the user")
	databaseName := cmd.String("database", "", "Database to grant permission on")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *username == "" || *permission == "" || *databaseName == "" {
		fmt.Println("Usage: dbac database grant-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.GrantPermissions(*databaseName, *username, *permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.GrantPermissions(*databaseName, *username, *permission)
		mysql.Close()
	}
}

func grantTable(params []string) {
	cmd := flag.NewFlagSet("grant-table", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the user")
	permission := cmd.String("permission", "", "Permission of the user")
	table := cmd.String("table", "", "Table to grant permission on")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *username == "" || *permission == "" || *table == "" {
		fmt.Println("Usage: dbac database grant-table --username [USERNAME] --permission [PERMISSION] --table [TABLE]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.GrantTablePermissions(*table, *username, *permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.GrantTablePermissions(currentProfile.Database, *table, *username, *permission)
		mysql.Close()
	}
}

func revokeDatabase(params []string) {
	cmd := flag.NewFlagSet("revoke-database", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the user")
	permission := cmd.String("permission", "", "Permission of the user")
	databaseName := cmd.String("database", "", "Database to revoke permission from")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *username == "" || *permission == "" || *databaseName == "" {
		fmt.Println("Usage: dbac database revoke-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.RevokePermissions(*databaseName, *username, *permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.RevokePermissions(*databaseName, *username, *permission)
		mysql.Close()
	}
}

func revokeTable(params []string) {
	cmd := flag.NewFlagSet("revoke-table", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the user")
	permission := cmd.String("permission", "", "Permission of the user")
	table := cmd.String("table", "", "Table to revoke permission from")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *username == "" || *permission == "" || *table == "" {
		fmt.Println("Usage: dbac database revoke-table --username [USERNAME] --permission [PERMISSION] --table [TABLE]")
		return
	}

	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.RevokeTablePermissions(*table, *username, *permission)
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.RevokeTablePermissions(currentProfile.Database, *table, *username, *permission)
		mysql.Close()
	}
}

func listDatabases() {
	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.ListDatabases()
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.ListDatabases()
		mysql.Close()
	}
}

func listTables() {
	switch currentProfile.DbType {
	case "psql":
		dbPort, _ := strconv.Atoi(currentProfile.Port)
		psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
		psql.ListTables()
		psql.Close()
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.ListTables()
		mysql.Close()
	}
}

func dumpDatabase(params []string) {
	cmd := flag.NewFlagSet("database-dump", flag.ExitOnError)
	path := cmd.String("path", "", "Query to be executed")
	database := cmd.String("database", "", "Query to be executed")
	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}
	if *database == "" {
		*database = currentProfile.Database
	}
	switch currentProfile.DbType {
	case "psql":
		fmt.Println("Not supported yet")
	case "mysql":
		mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
		mysql.Dump(*path, *database)
		mysql.Close()
	}

}

func execQuery(params []string) {
	cmd := flag.NewFlagSet("exec-query", flag.ExitOnError)
	query := cmd.String("query", "", "Query to be executed")
	file := cmd.String("file", "", "SQL file path")

	if err := cmd.Parse(params); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if *query == "" && *file == "" {
		fmt.Println("Usage: dbac database exec --query [QUERY] or --file [FILE]")
		return
	}

	if *file == "" {
		switch currentProfile.DbType {
		case "psql":
			dbPort, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.Exec(*query)
			psql.Close()
		case "mysql":
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.Exec(*query)
			mysql.Close()
		}
	} else {
		switch currentProfile.DbType {
		case "psql":
			dbPort, _ := strconv.Atoi(currentProfile.Port)
			psql.NewConnection(currentProfile.Host, dbPort, currentProfile.User, currentProfile.Password, currentProfile.Database)
			psql.FileExec(*file)
			psql.Close()
		case "mysql":
			mysql.NewConnection(currentProfile.Host, currentProfile.Port, currentProfile.User, currentProfile.Password, currentProfile.Database)
			mysql.FileExec(*file)
			mysql.Close()
		}
	}
}

func printDatabaseHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  dbac database ping")
	fmt.Println("  dbac database list-user")
	fmt.Println("  dbac database list-databases")
	fmt.Println("  dbac database list-tables")
	fmt.Println("  dbac database create-user --username [USERNAME] --user-password [PASSWORD]")
	fmt.Println("  dbac database create-database --database [DATABASE]")
	fmt.Println("  dbac database delete-user --username [USERNAME]")
	fmt.Println("  dbac database delete-database --database [DATABASE]")
	fmt.Println("  dbac database change-password --username [USERNAME] --new-password [NEW_PASSWORD]")
	fmt.Println("  dbac database grant-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE]")
	fmt.Println("  dbac database grant-table --username [USERNAME] --permission [PERMISSION] --table [TABLE]")
	fmt.Println("  dbac database revoke-database --username [USERNAME] --permission [PERMISSION] --database [DATABASE]")
	fmt.Println("  dbac database revoke-table --username [USERNAME] --permission [PERMISSION] --table [TABLE]")
	fmt.Println("  dbac database exec --query [QUERY] or --file [FILE]")
}
