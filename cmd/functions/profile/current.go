package functions

import (
	"dbac/cmd/helper"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CurrentProfileCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "current",
		Short:   "Show the current active profile",
		Long:    `Displays the name of the current active profile configured in the system.`,
		Example: "dbac profile current",
		Args:    cobra.NoArgs,
		Run:     currentProfiles,
	}
	subcommand.AddCommand(cmd)
}

func currentProfiles(cmd *cobra.Command, args []string) {
	currentProfileName, err := helper.GetCurrentProfileName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current profile name: %v\n", err)
		os.Exit(1)
	}

	profile := helper.ReadProfile(currentProfileName)
	fmt.Println("Current Profile Details:")
	fmt.Println("------------------------------")
	fmt.Printf("Name:        %s\n", profile.Name)
	fmt.Printf("Db Engine:   %s\n", profile.DbType)
	fmt.Printf("Host:        %s\n", profile.Host)
	fmt.Printf("Port:        %s\n", profile.Port)
	fmt.Printf("User:        %s\n", profile.User)
	// Commenting out password for security reasons, it shouldn't be printed
	// fmt.Printf("Password:    %s\n", profile.Password)
	fmt.Printf("Password:    %s\n", "getItOnProfileJson")
	fmt.Printf("Database:    %s\n", profile.Database)
	fmt.Println("------------------------------")
}
