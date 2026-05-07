package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadKeyRing_Missing(t *testing.T) {
	kr, err := LoadKeyRing("/nonexistent/path/keyring.json")
	if err != nil {
		t.Fatalf("expected nil error for missing keyring, got %v", err)
	}
	if len(kr.Entries) != 0 {
		t.Errorf("expected empty keyring, got %d entries", len(kr.Entries))
	}
}

func TestSaveAndLoadKeyRing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keyring.json")

	kr := &KeyRing{}
	kr.AddKey("alice", "pubkey-alice")
	kr.AddKey("bob", "pubkey-bob")

	if err := SaveKeyRing(path, kr); err != nil {
		t.Fatalf("SaveKeyRing: %v", err)
	}

	loaded, err := LoadKeyRing(path)
	if err != nil {
		t.Fatalf("LoadKeyRing: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded.Entries))
	}
}

func TestKeyRing_AddAndGet(t *testing.T) {
	kr := &KeyRing{}
	kr.AddKey("alice", "pubkey-alice")

	key, ok := kr.GetKey("alice")
	if !ok {
		t.Fatal("expected to find key 'alice'")
	}
	if key != "pubkey-alice" {
		t.Errorf("expected pubkey-alice, got %s", key)
	}
}

func TestKeyRing_AddKey_Overwrite(t *testing.T) {
	kr := &KeyRing{}
	kr.AddKey("alice", "old-key")
	kr.AddKey("alice", "new-key")

	if len(kr.Entries) != 1 {
		t.Errorf("expected 1 entry after overwrite, got %d", len(kr.Entries))
	}
	key, _ := kr.GetKey("alice")
	if key != "new-key" {
		t.Errorf("expected new-key, got %s", key)
	}
}

func TestKeyRing_RemoveKey(t *testing.T) {
	kr := &KeyRing{}
	kr.AddKey("alice", "pubkey-alice")

	removed := kr.RemoveKey("alice")
	if !removed {
		t.Fatal("expected RemoveKey to return true")
	}
	_, ok := kr.GetKey("alice")
	if ok {
		t.Error("expected key to be removed")
	}

	if kr.RemoveKey("ghost") {
		t.Error("expected false for missing key")
	}
}

func TestKeyRingFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keyring.json")

	kr := &KeyRing{}
	kr.AddKey("test", "pubkey")
	if err := SaveKeyRing(path, kr); err != nil {
		t.Fatalf("SaveKeyRing: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
