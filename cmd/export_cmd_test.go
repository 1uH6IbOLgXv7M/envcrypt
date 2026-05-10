package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envcrypt/crypto"
)

func setupExportTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	pubPath := filepath.Join(dir, "public.age")
	privPath := filepath.Join(dir, "private.age")
	if err := crypto.SaveKeyPair(pubPath, privPath, pub, priv); err != nil {
		t.Fatalf("SaveKeyPair: %v", err)
	}

	env := &crypto.EnvFile{Entries: []crypto.EnvEntry{{Key: "HELLO", Value: "world"}}}
	plain := []byte(env.Serialize())
	ct, err := crypto.Encrypt(pub, plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	vault := &crypto.Vault{}
	vault.Add(ct)
	if err := crypto.SaveVault(filepath.Join(dir, "vault.json"), vault); err != nil {
		t.Fatalf("SaveVault: %v", err)
	}
	return dir
}

func TestRunExport_NoArgs(t *testing.T) {
	err := runExport([]string{})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got: %v", err)
	}
}

func TestRunExport_DotEnv(t *testing.T) {
	dir := setupExportTestDir(t)
	vaultPath := filepath.Join(dir, "vault.json")
	outPath := filepath.Join(dir, "out.env")

	err := runExport([]string{vaultPath, "--format", "dotenv", "--out", outPath})
	if err != nil {
		t.Fatalf("runExport: %v", err)
	}
	b, _ := os.ReadFile(outPath)
	if !strings.Contains(string(b), "HELLO=world") {
		t.Errorf("expected HELLO=world in output, got: %s", b)
	}
}

func TestRunExport_JSON(t *testing.T) {
	dir := setupExportTestDir(t)
	vaultPath := filepath.Join(dir, "vault.json")
	outPath := filepath.Join(dir, "out.json")

	err := runExport([]string{vaultPath, "--format", "json", "--out", outPath})
	if err != nil {
		t.Fatalf("runExport json: %v", err)
	}
	b, _ := os.ReadFile(outPath)
	if !strings.Contains(string(b), `"HELLO"`) {
		t.Errorf("expected HELLO key in JSON, got: %s", b)
	}
}

func TestRunExport_Shell(t *testing.T) {
	dir := setupExportTestDir(t)
	vaultPath := filepath.Join(dir, "vault.json")
	outPath := filepath.Join(dir, "out.sh")

	err := runExport([]string{vaultPath, "--format", "shell", "--out", outPath})
	if err != nil {
		t.Fatalf("runExport shell: %v", err)
	}
	b, _ := os.ReadFile(outPath)
	if !strings.Contains(string(b), "export HELLO=") {
		t.Errorf("expected export HELLO= in output, got: %s", b)
	}
}

func TestRunExport_MissingVault(t *testing.T) {
	dir := t.TempDir()
	err := runExport([]string{filepath.Join(dir, "nope.json")})
	if err == nil {
		t.Error("expected error for missing vault")
	}
}
