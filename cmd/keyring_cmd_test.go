package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/crypto"
)

func setupKeyRingTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	_ = os.MkdirAll(".envcrypt", 0700)
	return dir, func() { _ = os.Chdir(origDir) }
}

func TestRunKeyRingAdd_NoArgs(t *testing.T) {
	if err := runKeyRingAdd([]string{}); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunKeyRingAdd_And_List(t *testing.T) {
	dir, cleanup := setupKeyRingTestDir(t)
	defer cleanup()

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	_ = priv

	pubFile := filepath.Join(dir, "alice.pub")
	if err := os.WriteFile(pubFile, []byte(pub), 0644); err != nil {
		t.Fatalf("writing pub key: %v", err)
	}

	if err := runKeyRingAdd([]string{"alice", pubFile}); err != nil {
		t.Fatalf("runKeyRingAdd: %v", err)
	}

	kr, err := crypto.LoadKeyRing(defaultKeyRingPath)
	if err != nil {
		t.Fatalf("LoadKeyRing: %v", err)
	}
	if len(kr.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(kr.Entries))
	}

	if err := runKeyRingList([]string{}); err != nil {
		t.Errorf("runKeyRingList: %v", err)
	}
}

func TestRunKeyRingRemove(t *testing.T) {
	dir, cleanup := setupKeyRingTestDir(t)
	defer cleanup()

	pub, _, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	pubFile := filepath.Join(dir, "bob.pub")
	_ = os.WriteFile(pubFile, []byte(pub), 0644)
	_ = runKeyRingAdd([]string{"bob", pubFile})

	if err := runKeyRingRemove([]string{"bob"}); err != nil {
		t.Fatalf("runKeyRingRemove: %v", err)
	}

	kr, _ := crypto.LoadKeyRing(defaultKeyRingPath)
	if len(kr.Entries) != 0 {
		t.Errorf("expected 0 entries after remove, got %d", len(kr.Entries))
	}
}

func TestRunKeyRingRemove_NoArgs(t *testing.T) {
	if err := runKeyRingRemove([]string{}); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunKeyRingList_Empty(t *testing.T) {
	_, cleanup := setupKeyRingTestDir(t)
	defer cleanup()

	if err := runKeyRingList([]string{}); err != nil {
		t.Errorf("runKeyRingList on empty: %v", err)
	}
}
