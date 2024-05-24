package functions

import (
	"dbac/cmd/helper"
	"fmt"

	_ "github.com/lib/pq"
)

var currentProfile helper.Profile

func App(params []string, profile helper.Profile) {
	currentProfile = profile

	if len(params) < 1 {
		printAppHelp()
		return
	}

	switch params[0] {
	case "batch":
		Batch(params)

	case "init":
		Init(params)

	case "profile":
		Profile(params)

	case "database":
		Database(params)

	case "-h":
		printAppHelp()

	default:
		fmt.Println("Invalid command.")
		printAppHelp()
	}
}

func printAppHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  dbac profile")
	fmt.Println("  dbac database")
	fmt.Println("  dbac batch")
	fmt.Println("  dbac init")
}
