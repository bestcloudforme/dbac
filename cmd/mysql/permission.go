package mysql

import "fmt"

func GrantPermissions(database, username, permissions string) error {
	perm, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	query := "GRANT " + perm + " ON " + quoteIdentifier(database) + ".* TO " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to grant permissions: %w", err)
	}
	fmt.Println("Permission added successfully")
	return nil
}

func RevokePermissions(database, username, permissions string) error {
	perm, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	query := "REVOKE " + perm + " ON " + quoteIdentifier(database) + ".* FROM " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to revoke permissions: %w", err)
	}
	fmt.Println("Permission revoked successfully")
	return nil
}

func GrantTablePermissions(database, table, username, permissions string) error {
	perm, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	query := "GRANT " + perm + " ON " + quoteIdentifier(database) + "." + quoteIdentifier(table) + " TO " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to grant table permissions: %w", err)
	}
	fmt.Println("Permission added successfully")
	return nil
}

func RevokeTablePermissions(database, table, username, permissions string) error {
	perm, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	query := "REVOKE " + perm + " ON " + quoteIdentifier(database) + "." + quoteIdentifier(table) + " FROM " + quoteLiteral(username) + "@'%';"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to revoke table permissions: %w", err)
	}
	fmt.Println("Permission revoked successfully")
	return nil
}
