package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	pg "github.com/habx/pg-commands"
)

func Dump2(path string, database string) error {
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

func Dump(path string, database string, filename string, table string, allTables bool) error {
	file, err := os.Create(path + filename)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "--\n-- dbac PostgreSQL database dump\n--\n")
	fmt.Fprintf(file, "SET statement_timeout = 0; -- Disables any timeout on statements to prevent long operations from being interrupted.\n")
	fmt.Fprintf(file, "SET lock_timeout = 0; -- Disables timeout on acquiring locks, ensuring that the restore process waits as long as needed.\n")
	fmt.Fprintf(file, "SET client_encoding = 'UTF8'; -- Sets the client encoding to UTF-8 to support a wide range of Unicode text.\n")
	fmt.Fprintf(file, "SET standard_conforming_strings = on; -- Enforces standard SQL behavior for interpreting string literals, avoiding misinterpretation of escape sequences.\n")
	fmt.Fprintf(file, "SELECT pg_catalog.set_config('search_path', '', false); -- Sets the search path to empty to avoid schema conflicts and ensure all references are fully qualified.\n")

	if table != "" {
		dumpTable(DbConnection, table, file)
	} else if allTables {
		dumpAllTables(DbConnection, file)
	} else {
		log.Fatal("No table specified and --all-tables flag is not set. Nothing to dump.")

	}

	return nil
}

func dumpTable(db *sql.DB, tableName string, file *os.File) {
	fmt.Fprintf(file, "--\n-- Name: %s; Type: TABLE; Schema: public; Owner: %s\n--\n", tableName, PostgresqlInfo.Username)
	dumpTableSchema(db, tableName, file)
	dumpTableData(db, tableName, file)
}

func dumpTableSchema(db *sql.DB, tableName string, file *os.File) {
	query := fmt.Sprintf("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name='%s'", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Fprintf(file, "CREATE TABLE public.%s (\n", tableName)
	first := true
	for rows.Next() {
		var colName, dataType, isNullable string
		if err := rows.Scan(&colName, &dataType, &isNullable); err != nil {
			log.Fatal(err)
		}
		if !first {
			fmt.Fprintf(file, ",\n")
		}
		first = false
		fmt.Fprintf(file, "    %s %s", colName, dataType)
		if isNullable == "NO" {
			fmt.Fprintf(file, " NOT NULL")
		}
	}
	fmt.Fprintf(file, "\n);\n\n")
	fmt.Fprintf(file, "ALTER TABLE public.%s OWNER TO %s;\n\n", tableName, PostgresqlInfo.Username)
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func dumpTableData(db *sql.DB, tableName string, file *os.File) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(file, "COPY public.%s (%s) FROM stdin;\n", tableName, strings.Join(columns, ", "))
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}

		var valueText []string
		for _, value := range values {
			if value == nil {
				valueText = append(valueText, "\\N")
			} else {
				colValue := escapeSQLString(string(value))
				if strings.ContainsAny(colValue, "\t\n\\") {
					colValue = fmt.Sprintf("\\%s", colValue)
				}
				valueText = append(valueText, colValue)
			}
		}
		fmt.Fprintf(file, "%s\n", strings.Join(valueText, "\t"))
	}
	fmt.Fprintf(file, "\\.\n\n")
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func dumpAllTables(db *sql.DB, file *os.File) {
	rows, err := db.Query("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		dumpTable(db, tableName, file)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func escapeSQLString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
