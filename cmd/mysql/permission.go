package mysql

import (
	"fmt"
	"log"
)

// GrantPermissions grants database-level permissions to a user
func GrantPermissions(database, username, permissions string) {
	query := fmt.Sprintf("GRANT %s ON %s.* TO '%s'@'%%';", permissions, database, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

// RevokePermissions revokes database-level permissions from a user
func RevokePermissions(database, username, permissions string) {
	query := fmt.Sprintf("REVOKE %s ON %s.* FROM '%s'@'%%';", permissions, database, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}

// GrantTablePermissions grants table-level permissions to a user
func GrantTablePermissions(database, table, username, permissions string) {
	query := fmt.Sprintf("GRANT %s ON %s.%s TO '%s'@'%%';", permissions, database, table, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

// RevokeTablePermissions revokes table-level permissions from a user
func RevokeTablePermissions(database, table, username, permissions string) {
	query := fmt.Sprintf("REVOKE %s ON %s.%s FROM '%s'@'%%';", permissions, database, table, username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}
