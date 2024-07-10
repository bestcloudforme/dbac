package functions

import (
	"dbac/cmd/helper"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func SwitchProfileCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "switch [profile-name]",
		Short:   "Switch to another profile",
		Long:    `Switches the current active profile to another one specified by name.`,
		Example: "dbac profile switch myprofile",
		Args:    cobra.ExactArgs(1),
		Run:     switchProfiles,
	}
	subcommand.AddCommand(cmd)
}

func switchProfiles(cmd *cobra.Command, args []string) {
	profileName := args[0]
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := os.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("Failed to read profile file: %v", err)
	}

	var allProfiles helper.Profiles
	if err := json.Unmarshal(data, &allProfiles); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	profileExists := false
	for _, profile := range allProfiles.Profiles {
		if profile.Name == profileName {
			profileExists = true
			break
		}
	}

	if !profileExists {
		fmt.Printf("Profile '%s' does not exist.\n", profileName)
		return
	}

	allProfiles.Current = profileName

	if err := helper.WriteProfilesToFile(allProfiles, profilePath); err != nil {
		log.Fatalf("Failed to write profile file: %v", err)
	}

	fmt.Printf("Switched to profile '%s' successfully.\n", profileName)
}
