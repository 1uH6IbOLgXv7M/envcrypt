package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"envcrypt/crypto"
)

// runRotate handles the `envcrypt rotate` command.
// Usage: envcrypt rotate <vault-file> <old-private-key-file> <new-public-key-file>
func runRotate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envcrypt rotate <vault-file> <old-private-key-file> <new-public-key-file>")
	}

	vaultPath := args[0]
	oldPrivKeyFile := args[1]
	newPubKeyFile := args[2]

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault file not found: %s", vaultPath)
	}

	oldPrivKey, err := crypto.LoadPrivateKey(oldPrivKeyFile)
	if err != nil {
		return fmt.Errorf("load old private key: %w", err)
	}

	newPubKey, err := crypto.LoadPublicKey(newPubKeyFile)
	if err != nil {
		return fmt.Errorf("load new public key: %w", err)
	}

	result, err := crypto.RotateVaultKeys(vaultPath, oldPrivKey, newPubKey)
	if err != nil {
		return fmt.Errorf("rotate keys: %w", err)
	}

	fmt.Printf("Key rotation complete: %d entr", result.EntriesRotated)
	if result.EntriesRotated == 1 {
		fmt.Print("y")
	} else {
		fmt.Print("ies")
	}
	fmt.Printf(" re-encrypted.\n")
	fmt.Printf("Vault: %s\n", filepath.Clean(vaultPath))
	return nil
}
