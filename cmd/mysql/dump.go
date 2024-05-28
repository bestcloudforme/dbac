package mysql

import (
	"fmt"

	"github.com/JamesStewy/go-mysqldump"
)

func Dump(path string, database string) error {
	dumpDir := path
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", database)
	dumper, err := mysqldump.Register(DbConnection, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Printf("Error registering database: %v\n", err)
		return err
	}

	defer dumper.Close()

	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Printf("Error dumping database: %v\n", err)
		return err
	}

	fmt.Printf("Dump successful, file saved to: %s\n", resultFilename)
	return nil
}
