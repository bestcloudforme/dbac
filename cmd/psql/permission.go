package psql

import (
	"fmt"
	"log"
)

func GrantPermissions(database, username, permissions string) {
	query := "GRANT " + validatePermissions(permissions) + " ON DATABASE " + quoteIdentifier(database) + " TO " + quoteIdentifier(username) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

func RevokePermissions(database, username, permissions string) {
	query := "REVOKE " + validatePermissions(permissions) + " ON DATABASE " + quoteIdentifier(database) + " FROM " + quoteIdentifier(username) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}

func GrantTablePermissions(table, username, permissions string) {
	query := "GRANT " + validatePermissions(permissions) + " ON " + quoteIdentifier(table) + " TO " + quoteIdentifier(username) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

func RevokeTablePermissions(table, username, permissions string) {
	query := "REVOKE " + validatePermissions(permissions) + " ON " + quoteIdentifier(table) + " FROM " + quoteIdentifier(username) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}
