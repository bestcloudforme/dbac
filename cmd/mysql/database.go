package mysql

import (
	"fmt"
	"log"
)

func ListDatabases() {
	query := "SHOW DATABASES;"
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
		log.Fatalf("Row iteration error: %v", err)
	}

	for _, db := range databases {
		fmt.Println(db)
	}
}

func ListTables() {
	query := "SHOW TABLES;"
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
		log.Fatalf("Row iteration error: %v", err)
	}

	for _, table := range tables {
		fmt.Println(table)
	}
}

func CreateDatabase(database string) {
	query := fmt.Sprintf("CREATE DATABASE %s;", database)
	_, err := DbConnection.Exec(query)
	if err != nil {
		log.Printf("Database couldn't be created: %v", err)
		return
	}
	fmt.Println("Database created successfully")
}

func DeleteDatabase(database string) {
	query := fmt.Sprintf("DROP DATABASE %s;", database)
	_, err := DbConnection.Exec(query)
	if err != nil {
		log.Printf("Database couldn't be deleted: %v", err)
		return
	}
	fmt.Println("Database deleted successfully")
}
