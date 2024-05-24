package main

import (
	"dbac/cmd/functions"
	"dbac/cmd/helper"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var currentProfile helper.Profile

func main() {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		handleFirstRun()
	} else {
		loadAndRun()
	}
}

func handleFirstRun() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	command := args[0]
	switch command {
	case "profile":
		if len(args) > 1 && args[1] == "add" {
			functions.App(args, currentProfile)
		} else {
			printUsage()
		}
	case "init":
		helper.CreateProfileFile()
	default:
		printUsage()
	}
}

func loadAndRun() {
	profile := helper.LoadProfile()
	currentProfile = helper.ReadProfile(profile)
	args := os.Args[1:]
	functions.App(args, currentProfile)
}

func printUsage() {
	fmt.Println("Please run one of these commands first:")
	fmt.Println("  dbac profile add --db-type [DB-TYPE] --host [HOST] --user [USER] --port [PORT] --password [PASSWORD] --database [DATABASE] --profile-name [NAME]")
	fmt.Println("  dbac init")
}
