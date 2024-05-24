package psql

import (
	"fmt"
	"log"
)

// GrantPermissions grants database-level permissions to a user
func GrantPermissions(database, username, permissions string) {
	query := fmt.Sprintf("GRANT %s ON DATABASE %s TO %s;", permissions, database, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

// RevokePermissions revokes database-level permissions from a user
func RevokePermissions(database, username, permissions string) {
	query := fmt.Sprintf("REVOKE %s ON DATABASE %s FROM %s;", permissions, database, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}

// GrantTablePermissions grants table-level permissions to a user
func GrantTablePermissions(table, username, permissions string) {
	query := fmt.Sprintf("GRANT %s ON %s TO %s;", permissions, table, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

// RevokeTablePermissions revokes table-level permissions from a user
func RevokeTablePermissions(table, username, permissions string) {
	query := fmt.Sprintf("REVOKE %s ON %s FROM %s;", permissions, table, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}
