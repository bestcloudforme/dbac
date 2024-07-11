package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var DbConnection *sql.DB

func NewConnection(host, port, user, password, dbname string) {
	addr := fmt.Sprintf("%s:%s", host, port)

	cfg := mysql.Config{
		User:    user,
		Passwd:  password,
		Net:     "tcp",
		Addr:    addr,
		DBName:  dbname,
		Timeout: 10 * time.Second,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	DbConnection = db
}

func Ping() {
	if err := DbConnection.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Connection Success")
}

func Close() {
	if err := DbConnection.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}

func Exec(query string) {
	rows, err := DbConnection.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		for i, col := range values {
			fmt.Printf("%s: %s\n", columns[i], col)
		}
		fmt.Println("---------------------")
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to iterate over rows: %v", err)
	}

	fmt.Println("Query run successfully")
}

func FileExec(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	commands := strings.Split(string(file), ";")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		if _, err := DbConnection.Exec(cmd); err != nil {
			log.Fatalf("Failed to execute command: %v", err)
		}
	}
	fmt.Println("SQL file run successfully")
}
