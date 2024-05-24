package psql

import (
	"fmt"
	"log"
)

// ListDatabases lists all databases in the PostgreSQL server
func ListDatabases() {
	query := "SELECT datname FROM pg_catalog.pg_database;"
	rows, err := DbConnection.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		databases = append(databases, database)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to iterate over rows: %v", err)
	}

	for _, database := range databases {
		fmt.Println(database)
	}
}

// ListTables lists all tables in the current PostgreSQL database
func ListTables() {
	query := "SELECT tablename FROM pg_catalog.pg_tables;"
	rows, err := DbConnection.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to iterate over rows: %v", err)
	}

	for _, table := range tables {
		fmt.Println(table)
	}
}

// CreateDatabase creates a new database in the PostgreSQL server
func CreateDatabase(database string) {
	query := fmt.Sprintf("CREATE DATABASE %s;", database)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Database couldn't be created: %v", err)
		return
	}
	fmt.Println("Database created successfully")
}

// DeleteDatabase deletes a database from the PostgreSQL server
func DeleteDatabase(database string) {
	query := fmt.Sprintf("DROP DATABASE %s;", database)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Database couldn't be deleted: %v", err)
		return
	}
	fmt.Println("Database deleted successfully")
}
