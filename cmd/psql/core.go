package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DbConnection *sql.DB
var PostgresqlInfo struct {
	Host     string
	Port     int
	DB       string
	Username string
	Password string
}

func NewConnection(host string, port int, user, password, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	PostgresqlInfo.DB = dbname
	PostgresqlInfo.Host = host
	PostgresqlInfo.Port = port
	PostgresqlInfo.Password = password
	PostgresqlInfo.Username = user

	DbConnection = db
}

func Ping() {
	if DbConnection == nil {
		log.Fatalf("No database connection established")
	}
	if err := DbConnection.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Connection Success")
}

func Close() {
	if DbConnection == nil {
		log.Fatalf("No database connection established")
	}
	if err := DbConnection.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}

func Exec(query string) {
	if DbConnection == nil {
		log.Fatalf("No database connection established")
	}
	if _, err := DbConnection.Exec(query); err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	fmt.Println("Query executed successfully")
}

func FileExec(filename string) {
	if DbConnection == nil {
		log.Fatalf("No database connection established")
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		request = strings.TrimSpace(request)
		if request == "" {
			continue
		}
		if _, err := DbConnection.Exec(request); err != nil {
			log.Fatalf("Failed to execute command: %v", err)
		}
	}
	fmt.Println("SQL file executed successfully")
}
