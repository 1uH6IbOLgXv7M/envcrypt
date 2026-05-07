package cmd

import (
	"fmt"
	"os"

	"github.com/user/envcrypt/crypto"
)

const defaultKeyRingPath = ".envcrypt/keyring.json"

// runKeyRingAdd adds a named public key to the keyring.
// Usage: envcrypt keyring add <name> <public-key-file>
func runKeyRingAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envcrypt keyring add <name> <public-key-file>")
	}
	name := args[0]
	keyFile := args[1]

	pubKey, err := crypto.LoadPublicKey(keyFile)
	if err != nil {
		return fmt.Errorf("loading public key: %w", err)
	}

	kr, err := crypto.LoadKeyRing(defaultKeyRingPath)
	if err != nil {
		return fmt.Errorf("loading keyring: %w", err)
	}

	if kr.HasKey(name) {
		return fmt.Errorf("key '%s' already exists in keyring; remove it first", name)
	}

	kr.AddKey(name, pubKey)

	if err := crypto.SaveKeyRing(defaultKeyRingPath, kr); err != nil {
		return fmt.Errorf("saving keyring: %w", err)
	}
	fmt.Printf("added key '%s' to keyring\n", name)
	return nil
}

// runKeyRingList prints all named keys in the keyring.
// Usage: envcrypt keyring list
func runKeyRingList(args []string) error {
	kr, err := crypto.LoadKeyRing(defaultKeyRingPath)
	if err != nil {
		return fmt.Errorf("loading keyring: %w", err)
	}
	if len(kr.Entries) == 0 {
		fmt.Println("keyring is empty")
		return nil
	}
	for _, e := range kr.Entries {
		fmt.Printf("%-20s %s  (added %s)\n", e.Name, e.PublicKey[:min(16, len(e.PublicKey))]+"...", e.CreatedAt.Format("2006-01-02"))
	}
	return nil
}

// runKeyRingRemove removes a named key from the keyring.
// Usage: envcrypt keyring remove <name>
func runKeyRingRemove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt keyring remove <name>")
	}
	name := args[0]

	kr, err := crypto.LoadKeyRing(defaultKeyRingPath)
	if err != nil {
		return fmt.Errorf("loading keyring: %w", err)
	}

	if !kr.RemoveKey(name) {
		fmt.Fprintf(os.Stderr, "key '%s' not found in keyring\n", name)
		return nil
	}

	if err := crypto.SaveKeyRing(defaultKeyRingPath, kr); err != nil {
		return fmt.Errorf("saving keyring: %w", err)
	}
	fmt.Printf("removed key '%s' from keyring\n", name)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
