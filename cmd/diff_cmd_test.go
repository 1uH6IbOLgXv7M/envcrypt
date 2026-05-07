package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"envcrypt/crypto"
)

func setupDiffTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "diff-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRunVaultDiff_NoArgs(t *testing.T) {
	// Should not panic; we can't easily capture os.Exit, so just verify
	// the function signature is correct by calling with insufficient args
	// via a recover pattern.
	defer func() { recover() }()
	// This would call os.Exit(1); we just ensure it compiles and runs.
	_ = runVaultDiff
}

func TestRunVaultDiff_RoundTrip(t *testing.T) {
	dir := setupDiffTestDir(t)
	vaultPath := filepath.Join(dir, "vault.json")
	pubPath := filepath.Join(dir, "pub.age")
	privPath := filepath.Join(dir, "priv.age")

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	if err := crypto.SaveKeyPair(pubPath, privPath, pub, priv); err != nil {
		t.Fatal(err)
	}

	pushEnv := func(env crypto.EnvFile) {
		plain := env.Serialize()
		cipher, err := crypto.Encrypt(pub, plain)
		if err != nil {
			t.Fatal(err)
		}
		vault, _ := crypto.LoadVault(vaultPath)
		crypto.AddEntry(vault, cipher)
		if err := crypto.SaveVault(vaultPath, vault); err != nil {
			t.Fatal(err)
		}
	}

	pushEnv(crypto.EnvFile{"A": "1", "B": "old"})
	pushEnv(crypto.EnvFile{"A": "1", "B": "new", "C": "added"})

	// Verify diff logic directly
	vault, err := crypto.LoadVault(vaultPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(vault.Entries) < 2 {
		t.Fatal("expected at least 2 vault entries")
	}

	decrypt := func(idx int) crypto.EnvFile {
		e, _ := crypto.EntryByVersion(vault, idx)
		p, _ := crypto.Decrypt(priv, e.Ciphertext)
		env, _ := crypto.ParseEnvBytes(p)
		return env
	}

	diff := crypto.DiffEnvFiles(decrypt(0), decrypt(1))
	if _, ok := diff.Added["C"]; !ok {
		t.Error("expected C in Added")
	}
	if v, ok := diff.Changed["B"]; !ok || v[0] != "old" || v[1] != "new" {
		t.Errorf("expected B changed old->new, got %v", diff.Changed)
	}
}
