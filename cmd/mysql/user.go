package mysql

import "fmt"

func ListUsers() error {
	rows, err := DbConnection.Query("SELECT user FROM mysql.user;")
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate over rows: %w", err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
	return nil
}

func CreateUser(username, password string) error {
	query := "CREATE USER " + quoteLiteral(username) + "@'%' IDENTIFIED BY " + quoteLiteral(password) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Println("User created successfully")
	return nil
}

func DeleteUser(username string) error {
	query := "DROP USER " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	fmt.Println("User deleted successfully")
	return nil
}

func ChangeUserPassword(username, newPassword string) error {
	query := "ALTER USER " + quoteLiteral(username) + "@'%' IDENTIFIED BY " + quoteLiteral(newPassword) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to change user password: %w", err)
	}
	fmt.Println("Password changed successfully")
	return nil
}
