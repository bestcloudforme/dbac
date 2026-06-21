package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
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
	dir := t.TempDir()
	setHome(t, dir)

	profiles := Profiles{
		Current: "prod",
		Profiles: []Profile{
			{Name: "dev", DbType: "psql", Host: "localhost", Port: "5432", User: "dev", Password: "secret", Database: "devdb"},
			{Name: "prod", DbType: "mysql", Host: "10.0.0.1", Port: "3306", User: "root", Password: "pass", Database: "proddb"},
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
