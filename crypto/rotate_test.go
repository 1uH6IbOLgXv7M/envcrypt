package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRotateVaultKeys_EmptyVault(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "vault.json")

	// Save an empty vault
	if err := SaveVault(vaultPath, &Vault{}); err != nil {
		t.Fatalf("save vault: %v", err)
	}

	_, newPub, _ := GenerateKeyPair()
	_, oldPriv, _ := GenerateKeyPair()

	result, err := RotateVaultKeys(vaultPath, oldPriv, newPub)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EntriesRotated != 0 {
		t.Errorf("expected 0 rotated, got %d", result.EntriesRotated)
	}
}

func TestRotateVaultKeys_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "vault.json")

	oldPub, oldPriv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate old key: %v", err)
	}
	newPub, newPriv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate new key: %v", err)
	}

	original := []byte("SECRET=hello\nOTHER=world")
	ciphertext, err := Encrypt(oldPub, original)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	vault := &Vault{}
	vault.AddEntry(ciphertext)
	if err := SaveVault(vaultPath, vault); err != nil {
		t.Fatalf("save vault: %v", err)
	}

	result, err := RotateVaultKeys(vaultPath, oldPriv, newPub)
	if err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if result.EntriesRotated != 1 {
		t.Errorf("expected 1 rotated, got %d", result.EntriesRotated)
	}

	updated, err := LoadVault(vaultPath)
	if err != nil {
		t.Fatalf("load updated vault: %v", err)
	}

	decrypted, err := Decrypt(newPriv, updated.Entries[0].Ciphertext)
	if err != nil {
		t.Fatalf("decrypt with new key: %v", err)
	}
	if string(decrypted) != string(original) {
		t.Errorf("expected %q, got %q", original, decrypted)
	}
}

func TestRotateVaultKeys_MissingVault(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "nonexistent.json")
	_, pub, _ := GenerateKeyPair()
	_, priv, _ := GenerateKeyPair()

	_, err := RotateVaultKeys(vaultPath, priv, pub)
	if err == nil {
		t.Error("expected error for missing vault, got nil")
	}
	_ = os.Remove(vaultPath)
}

func TestRotateVaultKeys_WrongOldKey(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "vault.json")

	oldPub, _, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate old key: %v", err)
	}
	newPub, _, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate new key: %v", err)
	}
	_, wrongPriv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate wrong key: %v", err)
	}

	ciphertext, _ := Encrypt(oldPub, []byte("KEY=val"))
	vault := &Vault{}
	vault.AddEntry(ciphertext)
	_ = SaveVault(vaultPath, vault)

	_, err = RotateVaultKeys(vaultPath, wrongPriv, newPub)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}
