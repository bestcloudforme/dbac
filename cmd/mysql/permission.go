package mysql

import (
	"fmt"
	"log"
)

func GrantPermissions(database, username, permissions string) {
	query := "GRANT " + validatePermissions(permissions) + " ON " + quoteIdentifier(database) + ".* TO " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

func RevokePermissions(database, username, permissions string) {
	query := "REVOKE " + validatePermissions(permissions) + " ON " + quoteIdentifier(database) + ".* FROM " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}

func GrantTablePermissions(database, table, username, permissions string) {
	query := "GRANT " + validatePermissions(permissions) + " ON " + quoteIdentifier(database) + "." + quoteIdentifier(table) + " TO " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be added: %v", err)
		return
	}
	fmt.Println("Permission added successfully")
}

func RevokeTablePermissions(database, table, username, permissions string) {
	query := "REVOKE " + validatePermissions(permissions) + " ON " + quoteIdentifier(database) + "." + quoteIdentifier(table) + " FROM " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Permission couldn't be revoked: %v", err)
		return
	}
	fmt.Println("Permission revoked successfully")
}
