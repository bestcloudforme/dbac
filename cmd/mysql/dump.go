package mysql

import (
	"fmt"

	"github.com/JamesStewy/go-mysqldump"
)

func Dump(path string, database string) {

	dumpDir := path
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", database)
	dumper, err := mysqldump.Register(DbConnection, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return
	}
	fmt.Printf("File is saved to %s", resultFilename)

	// Close dumper and connected database
	dumper.Close()
}
