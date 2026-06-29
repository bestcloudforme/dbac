package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var DbConnection *sql.DB

func NewConnection(host, port, user, password, dbname string) error {
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
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	DbConnection = db
	return nil
}

func Ping() error {
	if err := DbConnection.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Connection Success")
	return nil
}

func Close() error {
	if err := DbConnection.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

func Exec(query string) error {
	rows, err := DbConnection.Query(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		for i, col := range values {
			fmt.Printf("%s: %s\n", columns[i], col)
		}
		fmt.Println("---------------------")
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate over rows: %w", err)
	}
	fmt.Println("Query run successfully")
	return nil
}

func FileExec(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}
	commands := strings.Split(string(file), ";")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		if _, err := DbConnection.Exec(cmd); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
	}
	fmt.Println("SQL file run successfully")
	return nil
}

func quoteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func quoteLiteral(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return "'" + s + "'"
}

func validatePermissions(permissions string) (string, error) {
	for _, r := range permissions {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == ' ' || r == ',' || r == '_') {
			return "", fmt.Errorf("invalid permissions: %q", permissions)
		}
	}
	return permissions, nil
}
