package crypto

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func setupExportVault(t *testing.T) (vaultPath, pubPath, privPath string) {
	t.Helper()
	dir := t.TempDir()
	vaultPath = filepath.Join(dir, "vault.json")
	pubPath = filepath.Join(dir, "pub.age")
	privPath = filepath.Join(dir, "priv.age")

	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	if err := SaveKeyPair(pubPath, privPath, pub, priv); err != nil {
		t.Fatalf("SaveKeyPair: %v", err)
	}

	env := &EnvFile{Entries: []EnvEntry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}}
	plain := []byte(env.Serialize())
	ct, err := Encrypt(pub, plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	vault := &Vault{}
	vault.Add(ct)
	if err := SaveVault(vaultPath, vault); err != nil {
		t.Fatalf("SaveVault: %v", err)
	}
	return
}

func TestExportVault_Latest(t *testing.T) {
	vaultPath, _, privPath := setupExportVault(t)
	entry, err := ExportVault(vaultPath, privPath, 0)
	if err != nil {
		t.Fatalf("ExportVault: %v", err)
	}
	if entry.Version != 1 {
		t.Errorf("expected version 1, got %d", entry.Version)
	}
	if entry.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", entry.Vars["FOO"])
	}
}

func TestExportVault_MissingVault(t *testing.T) {
	dir := t.TempDir()
	_, err := ExportVault(filepath.Join(dir, "nope.json"), filepath.Join(dir, "priv.age"), 0)
	if err == nil {
		t.Error("expected error for missing vault")
	}
}

func TestRenderExport_JSON(t *testing.T) {
	entry := &ExportEntry{Version: 1, Timestamp: time.Now(), Vars: map[string]string{"A": "1"}}
	out, err := RenderExport(entry, FormatJSON)
	if err != nil {
		t.Fatalf("RenderExport JSON: %v", err)
	}
	if !strings.Contains(out, `"A"`) {
		t.Errorf("expected JSON to contain key A, got: %s", out)
	}
}

func TestRenderExport_DotEnv(t *testing.T) {
	entry := &ExportEntry{Version: 1, Vars: map[string]string{"X": "y"}}
	out, err := RenderExport(entry, FormatDotEnv)
	if err != nil {
		t.Fatalf("RenderExport dotenv: %v", err)
	}
	if !strings.Contains(out, "X=y") {
		t.Errorf("expected X=y in output, got: %s", out)
	}
}

func TestRenderExport_Shell(t *testing.T) {
	entry := &ExportEntry{Version: 1, Vars: map[string]string{"MY_VAR": "hello"}}
	out, err := RenderExport(entry, FormatShell)
	if err != nil {
		t.Fatalf("RenderExport shell: %v", err)
	}
	if !strings.Contains(out, "export MY_VAR=") {
		t.Errorf("expected export statement, got: %s", out)
	}
}

func TestRenderExport_UnknownFormat(t *testing.T) {
	entry := &ExportEntry{Vars: map[string]string{}}
	_, err := RenderExport(entry, ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWriteExport_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	if err := WriteExport("FOO=bar\n", path); err != nil {
		t.Fatalf("WriteExport: %v", err)
	}
	b, _ := os.ReadFile(path)
	if string(b) != "FOO=bar\n" {
		t.Errorf("unexpected file content: %s", b)
	}
}
