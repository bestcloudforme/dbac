package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	data, err := ioutil.ReadFile(profilePath)
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

func AddProfile(dbtype, host, user, password, database, port, name string) {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	newProfile := Profile{
		DbType:   dbtype,
		Host:     host,
		User:     user,
		Password: password,
		Database: database,
		Port:     port,
		Name:     name,
	}

	var allProfiles Profiles

	// Check if the profile file exists
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		allProfiles = Profiles{
			Profiles: []Profile{newProfile},
			Current:  name,
		}
	} else {
		data, err := ioutil.ReadFile(profilePath)
		if err != nil {
			log.Fatalf("Failed to read profile file: %v", err)
		}
		if len(data) == 0 {
			allProfiles = Profiles{
				Profiles: []Profile{},
				Current:  "",
			}
		} else {
			if err := json.Unmarshal(data, &allProfiles); err != nil {
				log.Fatalf("Failed to unmarshal profiles: %v", err)
			}
		}

		// Check if the profile name already exists
		for _, profile := range allProfiles.Profiles {
			if profile.Name == name {
				fmt.Printf("Profile name '%s' already exists.\n", name)
				return
			}
		}

		allProfiles.Profiles = append(allProfiles.Profiles, newProfile)
		allProfiles.Current = name
	}

	if err := writeProfilesToFile(allProfiles, profilePath); err != nil {
		log.Fatalf("Failed to write profile file: %v", err)
	}

	fmt.Println("Profile added successfully")
}

// Utility function to write profiles to file
func writeProfilesToFile(profiles Profiles, filePath string) error {
	jsonData, err := json.MarshalIndent(profiles, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, jsonData, 0644)
}

func AddFileProfile(file string) {
	var newProfile Profile
	var allProfiles Profiles

	yamlFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &newProfile); err != nil {
		log.Fatalf("Failed to unmarshal YAML file: %v", err)
	}

	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		allProfiles = Profiles{
			Profiles: []Profile{newProfile},
			Current:  newProfile.Name,
		}
	} else {
		data, err := os.ReadFile(profilePath)
		if err != nil {
			log.Fatalf("Failed to read profile file: %v", err)
		}
		if err := json.Unmarshal(data, &allProfiles); err != nil {
			log.Fatalf("Failed to unmarshal profile file: %v", err)
		}

		for _, profile := range allProfiles.Profiles {
			if profile.Name == newProfile.Name {
				fmt.Printf("Profile name '%s' already exists.\n", newProfile.Name)
				return
			}
		}

		allProfiles.Profiles = append(allProfiles.Profiles, newProfile)
		allProfiles.Current = newProfile.Name
	}

	if err := writeProfilesToFile(allProfiles, profilePath); err != nil {
		log.Fatalf("Failed to write profile file: %v", err)
	}

	fmt.Println("Profile added successfully")
}

func DeleteProfile(name string) {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	var allProfiles Profiles

	data, err := ioutil.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("Failed to read profile file: %v", err)
	}
	if err := json.Unmarshal(data, &allProfiles); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	// Filter out the profile to be deleted
	var newProfiles Profiles
	profileFound := false
	for _, profile := range allProfiles.Profiles {
		if profile.Name != name {
			newProfiles.Profiles = append(newProfiles.Profiles, profile)
		} else {
			profileFound = true
		}
	}

	if !profileFound {
		fmt.Printf("Profile '%s' not found.\n", name)
		return
	}

	// Update the current profile if the deleted profile was the current one
	if allProfiles.Current == name {
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
		if err := writeProfilesToFile(newProfiles, profilePath); err != nil {
			log.Fatalf("Failed to write profile file: %v", err)
		}
		fmt.Printf("Profile '%s' deleted successfully.\n", name)
	}
}

func ListProfiles() {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := ioutil.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No profiles found.")
			return
		}
		log.Fatalf("Failed to read profile file: %v", err)
	}

	var allProfiles Profiles
	if err := json.Unmarshal(data, &allProfiles); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	if len(allProfiles.Profiles) == 0 {
		fmt.Println("No profiles found.")
		return
	}

	for _, profile := range allProfiles.Profiles {
		if profile.Name == allProfiles.Current {
			fmt.Printf("* %s\n", profile.Name)
		} else {
			fmt.Printf("  %s\n", profile.Name)
		}
	}
}

func SwitchProfile(profilename string) {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := ioutil.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("Failed to read profile file: %v", err)
	}

	var allProfiles Profiles
	if err := json.Unmarshal(data, &allProfiles); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	// Check if the profile exists
	profileExists := false
	for _, profile := range allProfiles.Profiles {
		if profile.Name == profilename {
			profileExists = true
			break
		}
	}

	if !profileExists {
		fmt.Printf("Profile '%s' does not exist.\n", profilename)
		return
	}

	// Set the current profile
	allProfiles.Current = profilename

	// Write the updated profiles back to the file
	if err := writeProfilesToFile(allProfiles, profilePath); err != nil {
		log.Fatalf("Failed to write profile file: %v", err)
	}

	fmt.Printf("Switched to profile '%s' successfully.\n", profilename)
}

func CreateProfileFile() {
	profile := collectProfileData()

	var newProfiles Profiles
	newProfiles.Current = profile.Name
	newProfiles.Profiles = append(newProfiles.Profiles, profile)

	if err := writeProfilesToFile(newProfiles, os.Getenv("HOME")+"/.dbac-profiles.json"); err != nil {
		log.Fatalf("Failed to write profile file: %v", err)
	}
	fmt.Println("Profile created successfully")
}

func InitProfile() {
	profile := collectProfileData()
	var allProfiles Profiles

	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := ioutil.ReadFile(profilePath)
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
				fmt.Printf("Profile name '%s' already exists.\n", profile.Name)
				return
			}
		}

		allProfiles.Profiles = append(allProfiles.Profiles, profile)
	}

	if err := writeProfilesToFile(allProfiles, profilePath); err != nil {
		log.Fatalf("Failed to write profile file: %v", err)
	}

	fmt.Println("Profile initialized successfully")
}

func collectProfileData() Profile {
	var profile Profile
	var dbtype string
	var defaultPort string

	for {
		fmt.Println("Database Type (Choices: mysql, psql)")
		fmt.Print("-> ")
		fmt.Scan(&dbtype)
		if dbtype == "mysql" || dbtype == "psql" {
			profile.DbType = dbtype
			break
		}
		fmt.Println("Invalid choice. Please enter 'mysql' or 'psql'.")
	}

	if dbtype == "mysql" {
		defaultPort = "3306"
	} else if dbtype == "psql" {
		defaultPort = "5432"
	}

	fmt.Println("--------")
	fmt.Println("Database Host Address (e.g., localhost, 10.0.0.10, mydatabase.example.com)")
	fmt.Print("-> ")
	fmt.Scan(&profile.Host)

	fmt.Println("--------")
	fmt.Println("Database Username")
	fmt.Print("-> ")
	fmt.Scan(&profile.User)

	fmt.Println("--------")
	fmt.Println("Database Password")
	fmt.Print("-> ")
	fmt.Scan(&profile.Password)

	fmt.Println("--------")
	fmt.Printf("Database Name to Connect (e.g., postgres, mysql, testdb)\n-> ")
	fmt.Scan(&profile.Database)

	fmt.Println("--------")
	fmt.Printf("Database Port (Default: %s for %s)\n-> ", defaultPort, dbtype)
	fmt.Scan(&profile.Port)
	if profile.Port == "" {
		profile.Port = defaultPort
	}

	fmt.Println("--------")
	fmt.Println("Profile Name (how you will refer to this profile in CLI)")
	fmt.Print("-> ")
	fmt.Scan(&profile.Name)

	return profile
}
