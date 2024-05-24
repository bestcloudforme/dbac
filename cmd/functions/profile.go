package functions

import (
	"dbac/cmd/helper"
	"flag"
	"fmt"
)

func Profile(params []string) {
	if len(params) < 2 {
		printProfileHelp()
		return
	}

	switch params[1] {
	case "current":
		fmt.Println(currentProfile.Name)

	case "switch":
		if len(params) < 3 {
			fmt.Println("Usage: dbac profile switch [PROFILE-NAME]")
			return
		}
		helper.SwitchProfile(params[2])

	case "add":
		cmd := flag.NewFlagSet("add", flag.ExitOnError)
		dbType := cmd.String("db-type", "", "Type of the database (mysql or postgres)")
		host := cmd.String("host", "", "Database host")
		user := cmd.String("user", "", "Database user")
		port := cmd.String("port", "", "Database port")
		password := cmd.String("password", "", "Database password")
		database := cmd.String("database", "", "Database name")
		name := cmd.String("profile-name", "", "Profile name for pre-configured database connections")
		file := cmd.String("file", "", "Path to file with profile details")
		cmd.Parse(params[2:])

		if *file == "" {
			if *dbType == "" || *host == "" || *user == "" || *port == "" || *password == "" || *database == "" || *name == "" {
				fmt.Println("Usage: dbac profile add --db-type [DB-TYPE] --host [HOST] --user [USER] --port [PORT] --password [PASSWORD] --database [DATABASE] --profile-name [NAME]")
				return
			}
			helper.AddProfile(*dbType, *host, *user, *password, *database, *port, *name)
		} else {
			helper.AddFileProfile(*file)
		}

	case "delete":
		cmd := flag.NewFlagSet("delete", flag.ExitOnError)
		name := cmd.String("profile-name", "", "Profile name for pre-configured database connections")
		cmd.Parse(params[2:])
		if *name == "" {
			fmt.Println("Usage: dbac profile delete --profile-name [PROFILE-NAME]")
			return
		}
		helper.DeleteProfile(*name)

	case "list":
		helper.ListProfiles()

	case "-h":
		printProfileHelp()

	default:
		fmt.Println("Invalid command.")
		printProfileHelp()
	}
}

func printProfileHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  dbac profile current")
	fmt.Println("  dbac profile list")
	fmt.Println("  dbac profile switch [PROFILE-NAME]")
	fmt.Println("  dbac profile add --db-type [DB-TYPE] --host [HOST] --user [USER] --port [PORT] --password [PASSWORD] --database [DATABASE] --profile-name [NAME]")
	fmt.Println("  dbac profile add --file [FILE]")
	fmt.Println("  dbac profile delete --profile-name [PROFILE-NAME]")
}
