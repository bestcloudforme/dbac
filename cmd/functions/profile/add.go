package functions

import (
	"dbac/cmd/helper"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func AddProfileCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "add [--file filepath | [name] --db-type [dbtype] --user [username] --pass [password] --host [hostname] --port [port] --database [database]]",
		Short: "Add a new database profile",
		Long:  `Adds a new profile to the list of database connection profiles. You can add a profile by providing a YAML file or by specifying details manually.`,
		Example: `
# Add a profile by specifying details directly:
dbac profile add myprofile --db-type mysql --user admin --pass secure123 --host db.example.com --port 3306 --database sample_db

# Add a profile from a YAML file:
dbac profile add --file mysql-profile.yml
`,
		Args: cobra.MaximumNArgs(1),
		Run:  runAddProfile,
	}

	cmd.Flags().String("db-type", "", "Database type (mysql or postgres)")
	cmd.Flags().String("user", "", "Username for the database")
	cmd.Flags().String("pass", "", "Password for the database")
	cmd.Flags().String("host", "", "Host of the database")
	cmd.Flags().String("port", "", "Port of the database")
	cmd.Flags().String("database", "", "Name of the database")
	cmd.Flags().StringP("file", "f", "", "Path to the YAML file containing the profile details")

	subcommand.AddCommand(cmd)
}

func runAddProfile(cmd *cobra.Command, args []string) {
	file, _ := cmd.Flags().GetString("file")
	if file != "" {
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, "Cannot provide a profile name with --file flag\n")
			os.Exit(1)
		}
		profileName, err := helper.AddFileProfile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding profile from file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Profile \"%s\" added from file: %s\n", profileName, file)
	} else {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to display help: %v\n", err)
			}
			return
		}
		name := args[0]

		missingFlags := []string{}
		requiredFlags := []string{"db-type", "user", "pass", "host", "port", "database"}
		for _, flag := range requiredFlags {
			if !cmd.Flag(flag).Changed {
				missingFlags = append(missingFlags, flag)
			}
		}

		if len(missingFlags) > 0 {
			fmt.Fprintf(os.Stderr, "Missing required flags: %v\n", missingFlags)
			os.Exit(1)
		}

		dbtype, _ := cmd.Flags().GetString("db-type")
		username, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("pass")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		database, _ := cmd.Flags().GetString("database")

		profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
		newProfile := helper.Profile{
			DbType:   dbtype,
			Host:     host,
			User:     username,
			Password: password,
			Database: database,
			Port:     port,
			Name:     name,
		}

		var allProfiles helper.Profiles
		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			allProfiles = helper.Profiles{
				Profiles: []helper.Profile{newProfile},
				Current:  name,
			}
		} else {
			data, err := os.ReadFile(profilePath)
			if err != nil {
				log.Fatalf("Failed to read profile file: %v", err)
			}
			if len(data) == 0 {
				allProfiles = helper.Profiles{
					Profiles: []helper.Profile{},
					Current:  "",
				}
			} else {
				if err := json.Unmarshal(data, &allProfiles); err != nil {
					log.Fatalf("Failed to unmarshal profiles: %v", err)
				}
			}

			for _, profile := range allProfiles.Profiles {
				if profile.Name == name {
					fmt.Printf("Profile name '%s' already exists.\n", name)
					return
				}
			}

			allProfiles.Profiles = append(allProfiles.Profiles, newProfile)
			allProfiles.Current = name
		}

		if err := helper.WriteProfilesToFile(allProfiles, profilePath); err != nil {
			log.Fatalf("Failed to write profile file: %v", err)
		}

		fmt.Println("Profile added successfully")
	}
}
