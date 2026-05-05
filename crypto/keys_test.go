package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}

	dir := t.TempDir()
	if err := SaveKeyPair(kp, dir); err != nil {
		t.Fatalf("SaveKeyPair() error = %v", err)
	}

	pubPath := filepath.Join(dir, PublicKeyFile)
	privPath := filepath.Join(dir, PrivateKeyFile)

	// Check files exist
	if _, err := os.Stat(pubPath); err != nil {
		t.Errorf("public key file missing: %v", err)
	}
	if _, err := os.Stat(privPath); err != nil {
		t.Errorf("private key file missing: %v", err)
	}

	// Load and compare
	gotPub, err := LoadPublicKey(pubPath)
	if err != nil {
		t.Fatalf("LoadPublicKey() error = %v", err)
	}
	if gotPub != kp.PublicKey {
		t.Errorf("LoadPublicKey() = %q, want %q", gotPub, kp.PublicKey)
	}

	gotPriv, err := LoadPrivateKey(privPath)
	if err != nil {
		t.Fatalf("LoadPrivateKey() error = %v", err)
	}
	if gotPriv != kp.PrivateKey {
		t.Errorf("LoadPrivateKey() = %q, want %q", gotPriv, kp.PrivateKey)
	}
}

func TestSaveKeyPair_PrivateKeyPermissions(t *testing.T) {
	kp, _ := GenerateKeyPair()
	dir := t.TempDir()
	if err := SaveKeyPair(kp, dir); err != nil {
		t.Fatalf("SaveKeyPair() error = %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, PrivateKeyFile))
	if err != nil {
		t.Fatalf("stat private key: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("private key permissions = %o, want 0600", info.Mode().Perm())
	}
}

func TestLoadPublicKey_MissingFile(t *testing.T) {
	_, err := LoadPublicKey("/nonexistent/path/key.pub")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadKey_EmptyFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "key")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	_, err = LoadPublicKey(f.Name())
	if err == nil {
		t.Error("expected error for empty key file")
	}
}
