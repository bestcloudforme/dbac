package psql

import (
	"fmt"
	"os"

	pg "github.com/habx/pg-commands"
)

func Dump(path string, database string) {

	dump, err := pg.NewDump(&pg.Postgres{
		Host:     PostgresqlInfo.Host,
		Port:     PostgresqlInfo.Port,
		DB:       database,
		Username: PostgresqlInfo.Username,
		Password: PostgresqlInfo.Password,
	})
	if err != nil {
		panic(err)
	}
	a := dump.GetFileName()
	dump.SetFileName(path + "/" + a)

	dump.EnableVerbose()
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: true, StreamDestination: os.Stdout})
	if dumpExec.Error != nil {
		fmt.Println(dumpExec.Error.Err)
		fmt.Println(dumpExec.Output)

	} else {
		fmt.Println("Dump success")
		fmt.Println(dumpExec.File)
	}
}
