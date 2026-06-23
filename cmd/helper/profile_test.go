package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/zalando/go-keyring"
)

func writeProfileFile(t *testing.T, dir string, profiles Profiles) string {
	t.Helper()
	path := filepath.Join(dir, ".dbac-profiles.json")
	data, err := json.MarshalIndent(profiles, "", "    ")
	if err != nil {
		t.Fatalf("failed to marshal profiles: %v", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("failed to write profile file: %v", err)
	}
	return path
}

func setHome(t *testing.T, dir string) {
	t.Helper()
	original := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	t.Cleanup(func() { os.Setenv("HOME", original) })
}

func TestReadProfile(t *testing.T) {
	keyring.MockInit()

	dir := t.TempDir()
	setHome(t, dir)

	// Pre-store passwords in mock keychain; profiles on disk have no plaintext password.
	if err := StorePassword("dev", "secret"); err != nil {
		t.Fatalf("failed to seed keychain: %v", err)
	}
	if err := StorePassword("prod", "pass"); err != nil {
		t.Fatalf("failed to seed keychain: %v", err)
	}

	profiles := Profiles{
		Current: "prod",
		Profiles: []Profile{
			{Name: "dev", DbType: "psql", Host: "localhost", Port: "5432", User: "dev", Database: "devdb"},
			{Name: "prod", DbType: "mysql", Host: "10.0.0.1", Port: "3306", User: "root", Database: "proddb"},
		},
	}
	writeProfileFile(t, dir, profiles)

	t.Run("returns correct profile by name", func(t *testing.T) {
		p := ReadProfile("dev")
		if p.Name != "dev" {
			t.Errorf("expected name 'dev', got %q", p.Name)
		}
		if p.DbType != "psql" {
			t.Errorf("expected DbType 'psql', got %q", p.DbType)
		}
		if p.Port != "5432" {
			t.Errorf("expected port '5432', got %q", p.Port)
		}
		if p.Password != "secret" {
			t.Errorf("expected password 'secret' from keychain, got %q", p.Password)
		}
	})

	t.Run("psql profiles store 'psql' not 'postgres' as DbType", func(t *testing.T) {
		p := ReadProfile("dev")
		if p.DbType == "postgres" {
			t.Error("DbType must be 'psql', not 'postgres' — function switch cases use 'psql'")
		}
		if p.DbType != "psql" {
			t.Errorf("expected DbType 'psql', got %q", p.DbType)
		}
	})

	t.Run("returns second profile correctly", func(t *testing.T) {
		p := ReadProfile("prod")
		if p.DbType != "mysql" {
			t.Errorf("expected DbType 'mysql', got %q", p.DbType)
		}
		if p.Password != "pass" {
			t.Errorf("expected password 'pass' from keychain, got %q", p.Password)
		}
	})
}

func TestReadProfileMigration(t *testing.T) {
	keyring.MockInit()

	dir := t.TempDir()
	setHome(t, dir)

	// Profile has plaintext password in JSON — the legacy state before keychain.
	profiles := Profiles{
		Current: "legacy",
		Profiles: []Profile{
			{Name: "legacy", DbType: "mysql", Host: "localhost", Port: "3306", User: "admin", Password: "oldpass", Database: "mydb"},
		},
	}
	path := writeProfileFile(t, dir, profiles)

	p := ReadProfile("legacy")

	t.Run("returns correct password after migration", func(t *testing.T) {
		if p.Password != "oldpass" {
			t.Errorf("expected password 'oldpass', got %q", p.Password)
		}
	})

	t.Run("clears plaintext password from JSON after migration", func(t *testing.T) {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("could not read profile file: %v", err)
		}
		var got Profiles
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("could not unmarshal: %v", err)
		}
		if got.Profiles[0].Password != "" {
			t.Errorf("expected password cleared from JSON after migration, got %q", got.Profiles[0].Password)
		}
	})

	t.Run("password retrievable from keychain after migration", func(t *testing.T) {
		pw, err := RetrievePassword("legacy")
		if err != nil {
			t.Fatalf("failed to retrieve from keychain: %v", err)
		}
		if pw != "oldpass" {
			t.Errorf("expected 'oldpass' in keychain, got %q", pw)
		}
	})
}

func TestGetCurrentProfileName(t *testing.T) {
	dir := t.TempDir()
	setHome(t, dir)

	t.Run("returns current profile name", func(t *testing.T) {
		profiles := Profiles{
			Current:  "staging",
			Profiles: []Profile{{Name: "staging", DbType: "psql"}},
		}
		writeProfileFile(t, dir, profiles)

		name, err := GetCurrentProfileName()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "staging" {
			t.Errorf("expected 'staging', got %q", name)
		}
	})

	t.Run("returns error when current is empty", func(t *testing.T) {
		profiles := Profiles{
			Current:  "",
			Profiles: []Profile{{Name: "dev", DbType: "psql"}},
		}
		writeProfileFile(t, dir, profiles)

		_, err := GetCurrentProfileName()
		if err == nil {
			t.Error("expected error when current is empty, got nil")
		}
	})

	t.Run("returns error when profile file missing", func(t *testing.T) {
		emptyDir := t.TempDir()
		setHome(t, emptyDir)

		_, err := GetCurrentProfileName()
		if err == nil {
			t.Error("expected error when profile file does not exist, got nil")
		}
	})
}

func TestWriteProfilesToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")

	profiles := Profiles{
		Current: "local",
		Profiles: []Profile{
			{Name: "local", DbType: "psql", Host: "127.0.0.1", Port: "5432", User: "u", Password: "p", Database: "db"},
		},
	}

	if err := WriteProfilesToFile(profiles, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}

	var got Profiles
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("could not unmarshal written file: %v", err)
	}

	if got.Current != "local" {
		t.Errorf("expected current 'local', got %q", got.Current)
	}
	if len(got.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(got.Profiles))
	}
	if got.Profiles[0].DbType != "psql" {
		t.Errorf("expected DbType 'psql', got %q", got.Profiles[0].DbType)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("could not stat file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}
}
