package mysql

import (
	"fmt"
	"log"
)

func ListUsers() {
	query := "SELECT user FROM mysql.user;"
	rows, err := DbConnection.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to iterate over rows: %v", err)
	}

	for _, user := range users {
		fmt.Println(user)
	}
}

func CreateUser(username, password string) {
	query := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s';", username, password)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("User couldn't be created: %v", err)
		return
	}
	fmt.Println("User created successfully")
}

func DeleteUser(username string) {
	query := fmt.Sprintf("DROP USER '%s'@'%%';", username)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("User couldn't be deleted: %v", err)
		return
	}
	fmt.Println("User deleted successfully")
}

func ChangeUserPassword(username, newPassword string) {
	query := fmt.Sprintf("ALTER USER '%s'@'%%' IDENTIFIED BY '%s';", username, newPassword)
	if _, err := DbConnection.Exec(query); err != nil {
		log.Printf("Password couldn't be changed: %v", err)
		return
	}
	fmt.Println("Password changed successfully")
}
