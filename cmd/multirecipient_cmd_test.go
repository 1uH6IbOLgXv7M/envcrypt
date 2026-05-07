package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"envcrypt/crypto"
)

func setupMultiTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envcrypt-multi-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRunMultiEncrypt_NoArgs(t *testing.T) {
	if err := runMultiEncrypt([]string{}); err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRunMultiDecrypt_NoArgs(t *testing.T) {
	if err := runMultiDecrypt([]string{}); err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRunMultiEncryptDecrypt_RoundTrip(t *testing.T) {
	dir := setupMultiTestDir(t)

	pub1, priv1, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	recipFile := filepath.Join(dir, "recipients.json")
	t.Setenv("ENVCRYPT_RECIPIENTS_FILE", recipFile)

	rf, _ := crypto.LoadRecipientsFile(recipFile)
	rf.AddRecipient("alice", pub1)
	if err := crypto.SaveRecipientsFile(recipFile, rf); err != nil {
		t.Fatalf("save recipients: %v", err)
	}

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("KEY=value\nSECRET=abc"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	vaultFile := filepath.Join(dir, "vault.json")
	if err := runMultiEncrypt([]string{envFile, vaultFile}); err != nil {
		t.Fatalf("runMultiEncrypt: %v", err)
	}

	keyFile := filepath.Join(dir, "key.priv")
	if err := os.WriteFile(keyFile, []byte(priv1), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	outFile := filepath.Join(dir, "out.env")
	if err := runMultiDecrypt([]string{vaultFile, keyFile, outFile}); err != nil {
		t.Fatalf("runMultiDecrypt: %v", err)
	}

	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected non-empty decrypted output")
	}
}

func TestRunMultiEncrypt_NoRecipients(t *testing.T) {
	dir := setupMultiTestDir(t)
	recipFile := filepath.Join(dir, "recipients.json")
	t.Setenv("ENVCRYPT_RECIPIENTS_FILE", recipFile)

	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("K=V"), 0600)
	vaultFile := filepath.Join(dir, "vault.json")

	if err := runMultiEncrypt([]string{envFile, vaultFile}); err == nil {
		t.Fatal("expected error when no recipients")
	}
}
