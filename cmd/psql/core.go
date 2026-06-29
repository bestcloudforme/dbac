package psql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	pq "github.com/lib/pq"
)

var DbConnection *sql.DB

func NewConnection(host string, port int, user, password, dbname, sslmode string) error {
	if sslmode == "" {
		sslmode = "require"
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=10",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	DbConnection = db
	return nil
}

func Ping() error {
	if DbConnection == nil {
		return fmt.Errorf("no database connection established")
	}
	if err := DbConnection.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Connection Success")
	return nil
}

func Close() error {
	if DbConnection == nil {
		return fmt.Errorf("no database connection established")
	}
	if err := DbConnection.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

func Exec(query string) error {
	if DbConnection == nil {
		return fmt.Errorf("no database connection established")
	}
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	fmt.Println("Query executed successfully")
	return nil
}

func FileExec(filename string) error {
	if DbConnection == nil {
		return fmt.Errorf("no database connection established")
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		request = strings.TrimSpace(request)
		if request == "" {
			continue
		}
		if _, err := DbConnection.Exec(request); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
	}
	fmt.Println("SQL file executed successfully")
	return nil
}

func quoteIdentifier(name string) string { return pq.QuoteIdentifier(name) }
func quoteLiteral(s string) string       { return pq.QuoteLiteral(s) }

func validatePermissions(permissions string) (string, error) {
	for _, r := range permissions {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == ' ' || r == ',' || r == '_') {
			return "", fmt.Errorf("invalid permissions: %q", permissions)
		}
	}
	return permissions, nil
}
