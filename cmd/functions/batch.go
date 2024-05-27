package functions

import (
	"dbac/cmd/helper"
	"flag"
	"log"
)

func Batch(params []string) {
	cmd := flag.NewFlagSet("cmd", flag.ExitOnError)
	file := cmd.String("file", "", "File path for batch operations")
	if err := cmd.Parse(params[1:]); err != nil {
		log.Fatalf("Error parsing command line arguments: %v", err)
	}

	helper.RunBatch(*file)
}
