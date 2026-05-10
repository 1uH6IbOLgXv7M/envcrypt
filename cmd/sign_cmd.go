package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

	"envcrypt/crypto"
)

// runVaultSign signs the raw bytes of a vault file using an Ed25519 private key
// stored as a hex string in a key file, then prints the signature.
//
// Usage: envcrypt sign <vault-file> <signing-key-file>
func runVaultSign(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envcrypt sign <vault-file> <signing-key-file>")
	}
	vaultPath := args[0]
	keyPath := args[1]

	ciphertext, err := os.ReadFile(vaultPath)
	if err != nil {
		return fmt.Errorf("sign: cannot read vault file: %w", err)
	}

	keyHex, err := crypto.LoadPrivateKey(keyPath)
	if err != nil {
		return fmt.Errorf("sign: cannot load signing key: %w", err)
	}
	privBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("sign: invalid key encoding: %w", err)
	}

	sig, err := crypto.SignVault(privBytes, ciphertext)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}
	fmt.Printf("signature: %s\n", sig)
	return nil
}

// runVaultVerify verifies a vault file's signature using an Ed25519 public key.
//
// Usage: envcrypt verify <vault-file> <public-key-file> <signature-hex>
func runVaultVerify(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envcrypt verify <vault-file> <public-key-file> <signature-hex>")
	}
	vaultPath := args[0]
	keyPath := args[1]
	sigHex := args[2]

	ciphertext, err := os.ReadFile(vaultPath)
	if err != nil {
		return fmt.Errorf("verify: cannot read vault file: %w", err)
	}

	pubHex, err := crypto.LoadPublicKey(keyPath)
	if err != nil {
		return fmt.Errorf("verify: cannot load public key: %w", err)
	}
	pubBytes, err := hex.DecodeString(pubHex)
	if err != nil {
		return fmt.Errorf("verify: invalid key encoding: %w", err)
	}

	if err := crypto.VerifyVault(pubBytes, ciphertext, sigHex); err != nil {
		return fmt.Errorf("verify: %w", err)
	}
	fmt.Println("signature valid ✓")
	return nil
}
