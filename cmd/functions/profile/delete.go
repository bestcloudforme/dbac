package functions

import (
	"dbac/cmd/helper"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func DeleteProfileCommand(subcommand *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "delete [profile-name]",
		Short:   "Delete a profile",
		Long:    `Deletes a profile from the configuration by its name.`,
		Example: "dbac profile delete myprofile",
		Args:    cobra.ExactArgs(1),
		Run:     deleteProfiles,
	}
	subcommand.AddCommand(cmd)
}

func deleteProfiles(cmd *cobra.Command, args []string) {
	profileName := args[0]
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	var allProfiles helper.Profiles

	data, err := os.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("Failed to read profile file: %v", err)
	}
	if err := json.Unmarshal(data, &allProfiles); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	var newProfiles helper.Profiles
	profileFound := false
	for _, profile := range allProfiles.Profiles {
		if profile.Name != profileName {
			newProfiles.Profiles = append(newProfiles.Profiles, profile)
		} else {
			profileFound = true
		}
	}

	if !profileFound {
		fmt.Printf("Profile '%s' not found.\n", profileName)
		return
	}

	if allProfiles.Current == profileName {
		if len(newProfiles.Profiles) > 0 {
			newProfiles.Current = newProfiles.Profiles[0].Name
		} else {
			newProfiles.Current = ""
		}
	} else {
		newProfiles.Current = allProfiles.Current
	}

	if len(newProfiles.Profiles) == 0 {
		// Delete the profile file if no profiles remain
		if err := os.Remove(profilePath); err != nil {
			log.Fatalf("Failed to delete profile file: %v", err)
		}
		fmt.Println("All profiles deleted and profile file removed.")
	} else {
		if err := helper.WriteProfilesToFile(newProfiles, profilePath); err != nil {
			log.Fatalf("Failed to write profile file: %v", err)
		}
		fmt.Printf("Profile '%s' deleted successfully.\n", profileName)
	}
}
