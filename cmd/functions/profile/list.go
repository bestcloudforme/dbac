package functions

import (
	"dbac/cmd/helper"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func ListProfileCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all profiles",
		Long:    `Lists all configured database connection profiles.`,
		Example: "dbac profile list",
		Args:    cobra.NoArgs,
		Run:     listProfiles,
	}
	subcommand.AddCommand(cmd)
}

func listProfiles(cmd *cobra.Command, args []string) {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No profiles found.")
			return
		}
		log.Fatalf("Failed to read profile file: %v", err)
	}

	var allProfiles helper.Profiles
	if err := json.Unmarshal(data, &allProfiles); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	if len(allProfiles.Profiles) == 0 {
		fmt.Println("No profiles found.")
		return
	}

	fmt.Println("Configured Profiles:")
	for _, profile := range allProfiles.Profiles {
		if profile.Name == allProfiles.Current {
			fmt.Printf("* %s (current)\n", profile.Name) // Highlight the current profile
		} else {
			fmt.Printf("  %s\n", profile.Name)
		}
	}
}
