package functions

import (
	"dbac/cmd/helper"
	"os"
)

func Init(params []string) {
	if _, err := os.Stat(os.Getenv("HOME") + "/" + ".dbac-profiles.json"); err != nil {
		helper.CreateProfileFile()
	} else {
		helper.InitProfile()
	}
}
