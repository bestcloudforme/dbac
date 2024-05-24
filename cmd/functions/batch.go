package functions

import (
	"dbac/cmd/helper"
	"flag"
)

func Batch(params []string) {
	cmd := flag.NewFlagSet("cmd", flag.ExitOnError)
	file := cmd.String("file", "", "File path for batch operations")
	cmd.Parse(params[1:])

	helper.RunBatch(*file)
}
