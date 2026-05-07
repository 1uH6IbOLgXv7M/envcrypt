package cmd

import (
	"fmt"
	"os"

	"envcrypt/crypto"
)

// runMultiEncrypt encrypts a .env file for all recipients in the recipients file.
// Usage: envcrypt multi-encrypt <env-file> <vault-file>
func runMultiEncrypt(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envcrypt multi-encrypt <env-file> <vault-file>")
	}
	envFile, vaultFile := args[0], args[1]

	recipFile := recipientsFilePath()
	rf, err := crypto.LoadRecipientsFile(recipFile)
	if err != nil {
		return fmt.Errorf("load recipients: %w", err)
	}
	if len(rf.Recipients) == 0 {
		return fmt.Errorf("no recipients found in %s", recipFile)
	}

	envEntries, err := crypto.ParseEnvFile(envFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}
	plaintext := envEntries.Serialize()

	pubKeys := make(map[string]string, len(rf.Recipients))
	for alias, r := range rf.Recipients {
		pubKeys[alias] = r.PublicKey
	}

	ciphertexts, err := crypto.MultiRecipientEncrypt([]byte(plaintext), pubKeys)
	if err != nil {
		return fmt.Errorf("multi-encrypt: %w", err)
	}

	vault, _ := crypto.LoadVault(vaultFile)
	for alias, ct := range ciphertexts {
		vault.AddEntry(alias, ct)
	}
	if err := crypto.SaveVault(vaultFile, vault); err != nil {
		return fmt.Errorf("save vault: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Encrypted for %d recipients and saved to %s\n", len(pubKeys), vaultFile)
	return nil
}

// runMultiDecrypt decrypts a vault entry using the local private key.
// Usage: envcrypt multi-decrypt <vault-file> <private-key-file> <out-env-file>
func runMultiDecrypt(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envcrypt multi-decrypt <vault-file> <private-key-file> <out-env-file>")
	}
	vaultFile, keyFile, outFile := args[0], args[1], args[2]

	vault, err := crypto.LoadVault(vaultFile)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	privKey, err := crypto.LoadPrivateKey(keyFile)
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	latest := vault.LatestEntries()
	if len(latest) == 0 {
		return fmt.Errorf("vault is empty")
	}

	ciphertexts := make(map[string][]byte, len(latest))
	for alias, entry := range latest {
		ciphertexts[alias] = entry.Data
	}

	plaintext, alias, err := crypto.MultiRecipientDecrypt(ciphertexts, privKey)
	if err != nil {
		return fmt.Errorf("multi-decrypt: %w", err)
	}

	if err := os.WriteFile(outFile, plaintext, 0600); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Decrypted as recipient %q and wrote to %s\n", alias, outFile)
	return nil
}
