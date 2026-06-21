package functions

import (
	"testing"

	"github.com/spf13/cobra"
)

func newParent() *cobra.Command {
	return &cobra.Command{Use: "db"}
}

func findSubcommand(parent *cobra.Command, name string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func TestAddPingCommand(t *testing.T) {
	parent := newParent()
	AddPingCommand(parent)

	cmd := findSubcommand(parent, "ping")
	if cmd == nil {
		t.Fatal("ping command not registered")
	}
	if cmd.Args == nil {
		t.Error("expected Args validator to be set")
	}
}

func TestAddExecCommand(t *testing.T) {
	parent := newParent()
	AddExecCommand(parent)

	cmd := findSubcommand(parent, "exec")
	if cmd == nil {
		t.Fatal("exec command not registered")
	}

	for _, flag := range []string{"file", "query"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("exec command missing flag --%s", flag)
		}
	}

	if cmd.Flags().ShorthandLookup("f") == nil {
		t.Error("exec command missing shorthand -f for --file")
	}
}

func TestAddCreateDatabaseCommand(t *testing.T) {
	parent := newParent()
	AddCreateDatabaseCommand(parent)

	cmd := findSubcommand(parent, "create-database")
	if cmd == nil {
		t.Fatal("create-database command not registered")
	}
	if cmd.Flags().Lookup("database") == nil {
		t.Error("create-database command missing --database flag")
	}
}

func TestAddDeleteDatabaseCommand(t *testing.T) {
	parent := newParent()
	AddDeleteDatabaseCommand(parent)

	cmd := findSubcommand(parent, "delete-database")
	if cmd == nil {
		t.Fatal("delete-database command not registered")
	}
	if cmd.Flags().Lookup("database") == nil {
		t.Error("delete-database command missing --database flag")
	}
}

func TestAddCreateUserCommand(t *testing.T) {
	parent := newParent()
	AddCreateUserCommand(parent)

	cmd := findSubcommand(parent, "create-user")
	if cmd == nil {
		t.Fatal("create-user command not registered")
	}
	for _, flag := range []string{"username", "password"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("create-user command missing --%s flag", flag)
		}
	}
}

func TestAddDeleteUserCommand(t *testing.T) {
	parent := newParent()
	AddDeleteUserCommand(parent)

	cmd := findSubcommand(parent, "delete-user")
	if cmd == nil {
		t.Fatal("delete-user command not registered")
	}
	if cmd.Flags().Lookup("username") == nil {
		t.Error("delete-user command missing --username flag")
	}
}

func TestAddChangePasswordCommand(t *testing.T) {
	parent := newParent()
	AddChangePasswordCommand(parent)

	cmd := findSubcommand(parent, "change-password")
	if cmd == nil {
		t.Fatal("change-password command not registered")
	}
	for _, flag := range []string{"username", "new-password"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("change-password command missing --%s flag", flag)
		}
	}
}

func TestAddListDatabasesCommand(t *testing.T) {
	parent := newParent()
	AddListDatabasesCommand(parent)

	if findSubcommand(parent, "list-databases") == nil {
		t.Fatal("list-databases command not registered")
	}
}

func TestAddListTablesCommand(t *testing.T) {
	parent := newParent()
	AddListTablesCommand(parent)

	if findSubcommand(parent, "list-tables") == nil {
		t.Fatal("list-tables command not registered")
	}
}

func TestAddListUsersCommand(t *testing.T) {
	parent := newParent()
	AddListUsersCommand(parent)

	if findSubcommand(parent, "list-users") == nil {
		t.Fatal("list-users command not registered")
	}
}

func TestAddGrantDatabaseCommand(t *testing.T) {
	parent := newParent()
	AddGrantDatabaseCommand(parent)

	cmd := findSubcommand(parent, "grant-database")
	if cmd == nil {
		t.Fatal("grant-database command not registered")
	}
	for _, flag := range []string{"username", "permission", "database"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("grant-database command missing --%s flag", flag)
		}
	}
}

func TestAddGrantTableCommand(t *testing.T) {
	parent := newParent()
	AddGrantTableCommand(parent)

	cmd := findSubcommand(parent, "grant-table")
	if cmd == nil {
		t.Fatal("grant-table command not registered")
	}
	for _, flag := range []string{"username", "permission", "table"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("grant-table command missing --%s flag", flag)
		}
	}
}

func TestAddRevokeDatabaseCommand(t *testing.T) {
	parent := newParent()
	AddRevokeDatabaseCommand(parent)

	cmd := findSubcommand(parent, "revoke-database")
	if cmd == nil {
		t.Fatal("revoke-database command not registered")
	}
	for _, flag := range []string{"username", "permission", "database"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("revoke-database command missing --%s flag", flag)
		}
	}
}

func TestAddRevokeTableCommand(t *testing.T) {
	parent := newParent()
	AddRevokeTableCommand(parent)

	cmd := findSubcommand(parent, "revoke-table")
	if cmd == nil {
		t.Fatal("revoke-table command not registered")
	}
	for _, flag := range []string{"username", "permission", "table"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("revoke-table command missing --%s flag", flag)
		}
	}
}
