package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"envcrypt/crypto"
)

func setupVaultTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	if err := crypto.SaveKeyPair(filepath.Join(dir, ".env.pub"), filepath.Join(dir, ".env.key"), pub, priv); err != nil {
		t.Fatalf("SaveKeyPair: %v", err)
	}

	envContent := "APP_ENV=production\nSECRET=hunter2\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(envContent), 0600); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	return dir
}

func TestRunVaultPush(t *testing.T) {
	dir := setupVaultTestDir(t)
	vaultPath := filepath.Join(dir, ".env.vault")

	err := runVaultPush(filepath.Join(dir, ".env"), filepath.Join(dir, ".env.pub"), vaultPath)
	if err != nil {
		t.Fatalf("runVaultPush: %v", err)
	}

	v, err := crypto.LoadVault(vaultPath)
	if err != nil || len(v.Entries) != 1 {
		t.Fatalf("expected 1 vault entry, got %d (err: %v)", len(v.Entries), err)
	}
}

func TestRunVaultPushPull_RoundTrip(t *testing.T) {
	dir := setupVaultTestDir(t)
	vaultPath := filepath.Join(dir, ".env.vault")
	outPath := filepath.Join(dir, ".env.out")

	if err := runVaultPush(filepath.Join(dir, ".env"), filepath.Join(dir, ".env.pub"), vaultPath); err != nil {
		t.Fatalf("push: %v", err)
	}
	if err := runVaultPull(0, filepath.Join(dir, ".env.key"), vaultPath, outPath); err != nil {
		t.Fatalf("pull: %v", err)
	}

	original, _ := os.ReadFile(filepath.Join(dir, ".env"))
	restored, _ := os.ReadFile(outPath)
	if string(original) != string(restored) {
		t.Errorf("round-trip mismatch:\noriginal: %q\nrestored: %q", original, restored)
	}
}

func TestRunVaultPull_SpecificVersion(t *testing.T) {
	dir := setupVaultTestDir(t)
	vaultPath := filepath.Join(dir, ".env.vault")

	runVaultPush(filepath.Join(dir, ".env"), filepath.Join(dir, ".env.pub"), vaultPath)
	os.WriteFile(filepath.Join(dir, ".env"), []byte("APP_ENV=staging\n"), 0600)
	runVaultPush(filepath.Join(dir, ".env"), filepath.Join(dir, ".env.pub"), vaultPath)

	outPath := filepath.Join(dir, ".env.v1")
	if err := runVaultPull(1, filepath.Join(dir, ".env.key"), vaultPath, outPath); err != nil {
		t.Fatalf("pull version 1: %v", err)
	}
	data, _ := os.ReadFile(outPath)
	if string(data) != "APP_ENV=production\nSECRET=hunter2\n" {
		t.Errorf("unexpected v1 content: %q", data)
	}
}

func TestRunVaultPull_MissingVault(t *testing.T) {
	dir := setupVaultTestDir(t)
	err := runVaultPull(0, filepath.Join(dir, ".env.key"), filepath.Join(dir, "noexist.vault"), filepath.Join(dir, "out"))
	if err == nil {
		t.Error("expected error pulling from empty vault")
	}
}
