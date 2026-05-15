package crypto

import (
	"fmt"
	"os"
	"testing"
)

func setupSearchVault(t *testing.T, dir string) (pubKey, privKeyPath string) {
	t.Helper()
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	pubPath := dir + "/pub.age"
	privPath := dir + "/priv.age"
	if err := SaveKeyPair(pubPath, privPath, pub, priv); err != nil {
		t.Fatalf("SaveKeyPair: %v", err)
	}

	envLines := []EnvPair{
		{Key: "DATABASE_URL", Value: "postgres://localhost/mydb"},
		{Key: "SECRET_TOKEN", Value: "supersecret"},
		{Key: "PORT", Value: "8080"},
	}
	ef := &EnvFile{Pairs: envLines}
	plaintext := []byte(ef.Serialize())

	ct, err := Encrypt(pub, plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	vault := &Vault{}
	vault.AddEntry(ct)
	if err := SaveVault(dir, vault); err != nil {
		t.Fatalf("SaveVault: %v", err)
	}
	return pub, privPath
}

func TestSearchVault_MatchKey(t *testing.T) {
	dir := t.TempDir()
	_, privPath := setupSearchVault(t, dir)

	results, err := SearchVault(dir, privPath, "DATABASE")
	if err != nil {
		t.Fatalf("SearchVault: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", results[0].Key)
	}
}

func TestSearchVault_MatchValue(t *testing.T) {
	dir := t.TempDir()
	_, privPath := setupSearchVault(t, dir)

	results, err := SearchVault(dir, privPath, "supersecret")
	if err != nil {
		t.Fatalf("SearchVault: %v", err)
	}
	if len(results) != 1 || results[0].Key != "SECRET_TOKEN" {
		t.Errorf("unexpected results: %v", results)
	}
}

func TestSearchVault_NoMatch(t *testing.T) {
	dir := t.TempDir()
	_, privPath := setupSearchVault(t, dir)

	results, err := SearchVault(dir, privPath, "NONEXISTENT")
	if err != nil {
		t.Fatalf("SearchVault: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchVault_CaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	_, privPath := setupSearchVault(t, dir)

	results, err := SearchVault(dir, privPath, "port")
	if err != nil {
		t.Fatalf("SearchVault: %v", err)
	}
	if len(results) != 1 || results[0].Key != "PORT" {
		t.Errorf("expected PORT match, got %v", results)
	}
}

func TestSearchVault_MissingVault(t *testing.T) {
	dir := t.TempDir()
	_, privPath := dir+"/pub.age", dir+"/priv.age"
	_ = privPath

	// Write a dummy private key file so LoadPrivateKey doesn't fail first
	_, priv, _ := GenerateKeyPair()
	privFile := dir + "/priv.age"
	_ = os.WriteFile(privFile, []byte(fmt.Sprintf("%x", priv)), 0600)

	_, err := SearchVault(dir, privFile, "anything")
	if err == nil {
		t.Error("expected error for missing vault")
	}
}
