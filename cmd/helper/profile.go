package helper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Profile struct {
	DbType   string `json:"db-type" yaml:"db-type"`
	Host     string `json:"host" yaml:"host"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	Name     string `json:"name" yaml:"name"`
	Port     string `json:"port" yaml:"port"`
}

type Profiles struct {
	Profiles []Profile `json:"profiles"`
	Current  string    `json:"current"`
}

func ReadProfile(profileName string) Profile {
	var profiles Profiles
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := os.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("Failed to read the profile file: %v", err)
	}
	if len(data) == 0 {
		log.Fatalf("Profile file is empty")
	}
	if err := json.Unmarshal(data, &profiles); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}
	if len(profiles.Profiles) == 0 {
		log.Fatalf("No profiles found in the profile file")
	}
	for _, profile := range profiles.Profiles {
		if profile.Name == profileName {
			return profile
		}
	}

	log.Fatalf("Profile not found: %s", profileName)
	return Profile{}
}

func WriteProfilesToFile(profiles Profiles, filePath string) error {
	jsonData, err := json.MarshalIndent(profiles, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0644)
}

func AddFileProfile(file string) (string, error) {
	var newProfile Profile
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read YAML file: %v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &newProfile); err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML file: %v", err)
	}

	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	var allProfiles Profiles
	if data, err := os.ReadFile(profilePath); err == nil {
		if err := json.Unmarshal(data, &allProfiles); err != nil {
			return "", fmt.Errorf("failed to unmarshal profile file: %v", err)
		}
	}

	for _, profile := range allProfiles.Profiles {
		if profile.Name == newProfile.Name {
			return newProfile.Name, fmt.Errorf("profile name '%s' already exists", newProfile.Name)
		}
	}

	allProfiles.Profiles = append(allProfiles.Profiles, newProfile)
	allProfiles.Current = newProfile.Name

	if err := WriteProfilesToFile(allProfiles, profilePath); err != nil {
		return newProfile.Name, fmt.Errorf("failed to write profile file: %v", err)
	}

	return newProfile.Name, nil
}

func InitProfile() {
	profile := collectProfileData()
	var allProfiles Profiles

	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			allProfiles.Current = profile.Name
			allProfiles.Profiles = append(allProfiles.Profiles, profile)
		} else {
			log.Fatalf("Failed to read profile file: %v", err)
		}
	} else {
		if err := json.Unmarshal(data, &allProfiles); err != nil {
			log.Fatalf("Failed to unmarshal profiles: %v", err)
		}

		for _, p := range allProfiles.Profiles {
			if p.Name == profile.Name {
				log.Fatalf("Profile name '%s' already exists.\n", profile.Name)
				return
			}
		}

		allProfiles.Profiles = append(allProfiles.Profiles, profile)
	}

	if err := WriteProfilesToFile(allProfiles, profilePath); err != nil {
		log.Fatalf("Failed to write profile file: %v", err)
	}

	fmt.Println("Profile initialized successfully")
}

func collectProfileData() Profile {
	reader := bufio.NewReader(os.Stdin)
	var profile Profile
	var dbtype, defaultPort string
	var err error

	for {
		fmt.Println("[1/7] Database Type (Choices: mysql, psql):")
		fmt.Print("-> ")
		dbtype, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}
		dbtype = strings.TrimSpace(dbtype)

		if dbtype == "mysql" || dbtype == "psql" {
			profile.DbType = dbtype
			break
		} else {
			fmt.Println("Invalid choice. Please enter 'mysql' or 'psql'.")
		}
	}

	defaultPort = map[string]string{"mysql": "3306", "psql": "5432"}[dbtype]

	readRequiredField := func(prompt string) string {
		var input string
		for {
			fmt.Println(prompt)
			fmt.Print("-> ")
			input, err = reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Failed to read input: %v", err)
			}
			input = strings.TrimSpace(input)
			if input == "" {
				fmt.Println("This field cannot be empty. Please enter a valid value.")
			} else {
				break
			}
		}
		return input
	}

	profile.Host = readRequiredField("[2/7] Database Host Address (e.g., localhost, 10.0.0.10, mydatabase.example.com)")
	profile.User = readRequiredField("[3/7] Database Username")
	profile.Password = readRequiredField("[4/7] Database Password")
	profile.Database = readRequiredField("[5/7] Database Name to Connect (e.g., postgres, mysql, testdb)")

	fmt.Printf("[6/7] Database Port (Default: %s for %s):\n-> ", defaultPort, dbtype)
	if profile.Port, err = reader.ReadString('\n'); err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	profile.Port = strings.TrimSpace(profile.Port)
	if profile.Port == "" {
		profile.Port = defaultPort
	}
	profile.Name = readRequiredField("[7/7] Profile Name (how you will refer to this profile in CLI)")

	return profile
}

func GetCurrentProfileName() (string, error) {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read the profile file: %v", err)
	}

	var profiles Profiles
	if err := json.Unmarshal(data, &profiles); err != nil {
		return "", fmt.Errorf("failed to unmarshal profiles: %v", err)
	}

	if profiles.Current == "" {
		return "", fmt.Errorf("no current profile set")
	}

	return profiles.Current, nil
}
