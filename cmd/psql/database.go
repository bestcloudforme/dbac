package psql

import "fmt"

func ListDatabases() error {
	rows, err := DbConnection.Query("SELECT datname FROM pg_catalog.pg_database;")
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate over rows: %w", err)
	}
	for _, database := range databases {
		fmt.Println(database)
	}
	return nil
}

func ListTables() error {
	rows, err := DbConnection.Query("SELECT tablename FROM pg_catalog.pg_tables;")
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate over rows: %w", err)
	}
	for _, table := range tables {
		fmt.Println(table)
	}
	return nil
}

func CreateDatabase(database string) error {
	if _, err := DbConnection.Exec("CREATE DATABASE " + quoteIdentifier(database) + ";"); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	fmt.Println("Database created successfully")
	return nil
}

func DeleteDatabase(database string) error {
	if _, err := DbConnection.Exec("DROP DATABASE " + quoteIdentifier(database) + ";"); err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}
	fmt.Println("Database deleted successfully")
	return nil
}
