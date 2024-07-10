package helper

import (
	"encoding/json"
	"log"
	"os"
)

func LoadProfile() string {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := os.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("Failed to read profile file: %v", err)
	}

	var curProfile Profiles
	if err := json.Unmarshal(data, &curProfile); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	return curProfile.Current
}
