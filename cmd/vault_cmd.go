package cmd

import (
	"fmt"
	"os"

	"envcrypt/crypto"
)

const (
	defaultVaultPath  = ".env.vault"
	defaultEnvPath    = ".env"
	defaultPublicKey  = ".env.pub"
	defaultPrivateKey = ".env.key"
)

// runVaultPush encrypts the current .env and appends it to the vault.
func runVaultPush(envPath, pubKeyPath, vaultPath string) error {
	pubKey, err := crypto.LoadPublicKey(pubKeyPath)
	if err != nil {
		return fmt.Errorf("load public key: %w", err)
	}

	envFile, err := crypto.ParseEnvFile(envPath)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}

	plaintext := []byte(envFile.Serialize())
	ciphertext, err := crypto.Encrypt(pubKey, plaintext)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	vault, err := crypto.LoadVault(vaultPath)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	entry := vault.AddEntry(ciphertext)
	if err := crypto.SaveVault(vaultPath, vault); err != nil {
		return fmt.Errorf("save vault: %w", err)
	}

	fmt.Fprintf(os.Stdout, "pushed version %d to %s\n", entry.Version, vaultPath)
	return nil
}

// runVaultPull decrypts the latest (or specified) vault version to .env.
func runVaultPull(version int, privKeyPath, vaultPath, outPath string) error {
	privKey, err := crypto.LoadPrivateKey(privKeyPath)
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	vault, err := crypto.LoadVault(vaultPath)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	var entry crypto.VaultEntry
	if version <= 0 {
		entry, err = vault.LatestEntry()
	} else {
		entry, err = vault.EntryByVersion(version)
	}
	if err != nil {
		return fmt.Errorf("get vault entry: %w", err)
	}

	plaintext, err := crypto.Decrypt(privKey, entry.Ciphertext)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	if err := os.WriteFile(outPath, plaintext, 0600); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "pulled version %d to %s\n", entry.Version, outPath)
	return nil
}
