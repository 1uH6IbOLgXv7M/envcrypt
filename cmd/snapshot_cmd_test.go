package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"envcrypt/crypto"
)

func setupSnapshotTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func TestRunSnapshotList_Empty(t *testing.T) {
	dir := setupSnapshotTestDir(t)
	// Capture stdout via redirect trick using a pipe
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSnapshotList([]string{dir})

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No snapshots found") {
		t.Errorf("expected 'No snapshots found', got: %s", buf.String())
	}
}

func TestRunSnapshotList_WithEntries(t *testing.T) {
	dir := setupSnapshotTestDir(t)
	for i := 1; i <= 2; i++ {
		snap := crypto.Snapshot{
			Version:   i,
			Timestamp: time.Now().UTC(),
			Keys:      []string{"A", "B"},
			Checksum:  fmt.Sprintf("chk%d", i),
		}
		if err := crypto.AppendSnapshot(dir, snap); err != nil {
			t.Fatalf("append: %v", err)
		}
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSnapshotList([]string{dir})

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "VERSION") {
		t.Errorf("expected header, got: %s", output)
	}
	if !strings.Contains(output, "chk1") {
		t.Errorf("expected chk1 in output, got: %s", output)
	}
}

func TestRunSnapshotCreate_NoArgs(t *testing.T) {
	err := runSnapshotCreate([]string{})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRunSnapshotCreate_RoundTrip(t *testing.T) {
	dir := setupSnapshotTestDir(t)

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	pubFile := filepath.Join(dir, "pub.age")
	privFile := filepath.Join(dir, "priv.age")
	if err := crypto.SaveKeyPair(pubFile, privFile, pub, priv); err != nil {
		t.Fatalf("save keys: %v", err)
	}

	envContent := []byte("DB_URL=postgres://localhost\nAPI_KEY=secret\n")
	cipher, err := crypto.Encrypt(pub, envContent)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	vault := &crypto.Vault{}
	crypto.AddEntry(vault, cipher)
	if err := crypto.SaveVault(dir, vault); err != nil {
		t.Fatalf("save vault: %v", err)
	}

	err = runSnapshotCreate([]string{dir, privFile})
	if err != nil {
		t.Fatalf("snapshot create: %v", err)
	}

	log, err := crypto.LoadSnapshotLog(dir)
	if err != nil {
		t.Fatalf("load log: %v", err)
	}
	if len(log.Snapshots) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(log.Snapshots))
	}
	if len(log.Snapshots[0].Keys) == 0 {
		t.Error("expected keys in snapshot")
	}
}

func TestRunSnapshotList_DefaultDir(t *testing.T) {
	// Should not error even with no args (uses ".")
	// We just ensure it doesn't panic on a valid directory
	old, _ := os.Getwd()
	tmp := t.TempDir()
	os.Chdir(tmp)
	defer os.Chdir(old)

	err := runSnapshotList([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
