package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadVault_Missing(t *testing.T) {
	v, err := LoadVault("/tmp/nonexistent_vault_xyz.json")
	if err != nil {
		t.Fatalf("expected no error for missing vault, got: %v", err)
	}
	if len(v.Entries) != 0 {
		t.Fatalf("expected empty vault, got %d entries", len(v.Entries))
	}
}

func TestSaveAndLoadVault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault.json")

	v := &Vault{}
	v.AddEntry([]byte("ciphertext-one"))
	v.AddEntry([]byte("ciphertext-two"))

	if err := SaveVault(path, v); err != nil {
		t.Fatalf("SaveVault: %v", err)
	}

	loaded, err := LoadVault(path)
	if err != nil {
		t.Fatalf("LoadVault: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Version != 1 || loaded.Entries[1].Version != 2 {
		t.Errorf("unexpected versions: %v %v", loaded.Entries[0].Version, loaded.Entries[1].Version)
	}
}

func TestAddEntry_Versioning(t *testing.T) {
	v := &Vault{}
	e1 := v.AddEntry([]byte("a"))
	e2 := v.AddEntry([]byte("b"))
	e3 := v.AddEntry([]byte("c"))
	if e1.Version != 1 || e2.Version != 2 || e3.Version != 3 {
		t.Errorf("versions not sequential: %d %d %d", e1.Version, e2.Version, e3.Version)
	}
}

func TestLatestEntry(t *testing.T) {
	v := &Vault{}
	_, err := v.LatestEntry()
	if err == nil {
		t.Fatal("expected error on empty vault")
	}
	v.AddEntry([]byte("first"))
	v.AddEntry([]byte("second"))
	entry, err := v.LatestEntry()
	if err != nil {
		t.Fatalf("LatestEntry: %v", err)
	}
	if entry.Version != 2 {
		t.Errorf("expected version 2, got %d", entry.Version)
	}
}

func TestEntryByVersion(t *testing.T) {
	v := &Vault{}
	v.AddEntry([]byte("v1"))
	v.AddEntry([]byte("v2"))

	e, err := v.EntryByVersion(1)
	if err != nil || string(e.Ciphertext) != "v1" {
		t.Errorf("EntryByVersion(1) failed: %v", err)
	}
	_, err = v.EntryByVersion(99)
	if err == nil {
		t.Error("expected error for missing version")
	}
}

func TestLoadVault_Corrupt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")
	os.WriteFile(path, []byte("not valid json{"), 0644)
	_, err := LoadVault(path)
	if err == nil {
		t.Fatal("expected error for corrupt vault")
	}
}

func TestSaveAndLoadVault_CiphertextRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "roundtrip.vault.json")

	original := []byte("super-secret-ciphertext")
	v := &Vault{}
	v.AddEntry(original)

	if err := SaveVault(path, v); err != nil {
		t.Fatalf("SaveVault: %v", err)
	}

	loaded, err := LoadVault(path)
	if err != nil {
		t.Fatalf("LoadVault: %v", err)
	}

	entry, err := loaded.EntryByVersion(1)
	if err != nil {
		t.Fatalf("EntryByVersion: %v", err)
	}
	if string(entry.Ciphertext) != string(original) {
		t.Errorf("ciphertext mismatch: got %q, want %q", entry.Ciphertext, original)
	}
}
