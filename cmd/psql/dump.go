package psql

import (
	"fmt"
	"os"

	pg "github.com/habx/pg-commands"
)

func Dump(path string, database string) error {
	dump, err := pg.NewDump(&pg.Postgres{
		Host:     PostgresqlInfo.Host,
		Port:     PostgresqlInfo.Port,
		DB:       database,
		Username: PostgresqlInfo.Username,
		Password: PostgresqlInfo.Password,
	})
	if err != nil {
		fmt.Printf("Error creating dump object: %v\n", err)
		return err
	}

	fileName := dump.GetFileName()
	dump.SetFileName(path + "/" + fileName)

	dump.EnableVerbose()
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: true, StreamDestination: os.Stdout})
	if dumpExec.Error != nil {
		fmt.Printf("Error during dump execution: %v\nOutput: %s\n", dumpExec.Error.Err, dumpExec.Output)
		return dumpExec.Error.Err
	}

	fmt.Printf("Dump successful, file saved to: %s\n", dumpExec.File)
	return nil
}
