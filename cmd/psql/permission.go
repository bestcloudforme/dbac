package psql

import "fmt"

func GrantPermissions(database, username, permissions string) error {
	perm, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	query := "GRANT " + perm + " ON DATABASE " + quoteIdentifier(database) + " TO " + quoteIdentifier(username) + ";"
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
	query := "REVOKE " + perm + " ON DATABASE " + quoteIdentifier(database) + " FROM " + quoteIdentifier(username) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to revoke permissions: %w", err)
	}
	fmt.Println("Permission revoked successfully")
	return nil
}

func GrantTablePermissions(table, username, permissions string) error {
	perm, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	query := "GRANT " + perm + " ON " + quoteIdentifier(table) + " TO " + quoteIdentifier(username) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to grant table permissions: %w", err)
	}
	fmt.Println("Permission added successfully")
	return nil
}

func RevokeTablePermissions(table, username, permissions string) error {
	perm, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	query := "REVOKE " + perm + " ON " + quoteIdentifier(table) + " FROM " + quoteIdentifier(username) + ";"
	if _, err := DbConnection.Exec(query); err != nil {
		return fmt.Errorf("failed to revoke table permissions: %w", err)
	}
	fmt.Println("Permission revoked successfully")
	return nil
}
