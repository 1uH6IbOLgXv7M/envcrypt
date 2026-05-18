package crypto

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadProfileStore_Missing(t *testing.T) {
	dir := t.TempDir()
	store, err := LoadProfileStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Profiles) != 0 {
		t.Errorf("expected empty store, got %d profiles", len(store.Profiles))
	}
}

func TestSaveAndLoadProfileStore(t *testing.T) {
	dir := t.TempDir()
	store := &ProfileStore{Profiles: make(map[string]Profile)}
	store.AddProfile(Profile{Name: "dev", EnvFile: ".env.dev"})
	if err := SaveProfileStore(dir, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadProfileStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, ok := loaded.Profiles["dev"]; !ok {
		t.Error("expected profile 'dev' to exist")
	}
}

func TestProfileStore_AddAndGet(t *testing.T) {
	store := &ProfileStore{Profiles: make(map[string]Profile)}
	store.AddProfile(Profile{Name: "prod", EnvFile: ".env.prod", Metadata: map[string]string{"region": "us-east-1"}})
	p, ok := store.GetProfile("prod")
	if !ok {
		t.Fatal("expected profile to exist")
	}
	if p.EnvFile != ".env.prod" {
		t.Errorf("expected .env.prod, got %s", p.EnvFile)
	}
	if p.Metadata["region"] != "us-east-1" {
		t.Errorf("unexpected metadata: %v", p.Metadata)
	}
}

func TestProfileStore_AddProfile_Overwrite(t *testing.T) {
	store := &ProfileStore{Profiles: make(map[string]Profile)}
	store.AddProfile(Profile{Name: "staging", EnvFile: ".env.staging"})
	origCreated := store.Profiles["staging"].CreatedAt
	time.Sleep(2 * time.Millisecond)
	store.AddProfile(Profile{Name: "staging", EnvFile: ".env.staging.v2"})
	p, _ := store.GetProfile("staging")
	if p.EnvFile != ".env.staging.v2" {
		t.Errorf("expected updated env file")
	}
	if !p.CreatedAt.Equal(origCreated) {
		t.Errorf("CreatedAt should be preserved on overwrite")
	}
	if !p.UpdatedAt.After(origCreated) {
		t.Errorf("UpdatedAt should be after CreatedAt")
	}
}

func TestProfileStore_RemoveProfile(t *testing.T) {
	store := &ProfileStore{Profiles: make(map[string]Profile)}
	store.AddProfile(Profile{Name: "test", EnvFile: ".env.test"})
	if !store.RemoveProfile("test") {
		t.Error("expected remove to return true")
	}
	if store.RemoveProfile("test") {
		t.Error("expected remove of missing profile to return false")
	}
}

func TestProfileStorePath(t *testing.T) {
	dir := "/tmp/myproject"
	expected := filepath.Join(dir, ".envcrypt", "profiles.json")
	if got := ProfileStorePath(dir); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestSaveProfileStore_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	store := &ProfileStore{Profiles: make(map[string]Profile)}
	store.AddProfile(Profile{Name: "ci", EnvFile: ".env.ci"})
	if err := SaveProfileStore(dir, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(ProfileStorePath(dir))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
