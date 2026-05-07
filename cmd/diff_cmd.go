package cmd

import (
	"fmt"
	"os"
	"strconv"

	"envcrypt/crypto"
)

// runVaultDiff compares two vault versions and prints a diff.
// Usage: envcrypt diff <vault> <privkey> [versionA] [versionB]
// If versions are omitted, compares the last two entries.
func runVaultDiff(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envcrypt diff <vault> <privkey> [vA] [vB]")
		os.Exit(1)
	}
	vaultPath := args[0]
	privKeyPath := args[1]

	vault, err := crypto.LoadVault(vaultPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load vault: %v\n", err)
		os.Exit(1)
	}

	privKey, err := crypto.LoadPrivateKey(privKeyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load private key: %v\n", err)
		os.Exit(1)
	}

	if len(vault.Entries) < 2 {
		fmt.Println("vault has fewer than 2 entries; nothing to diff")
		return
	}

	vA := len(vault.Entries) - 2
	vB := len(vault.Entries) - 1

	if len(args) >= 4 {
		vA, _ = strconv.Atoi(args[2])
		vB, _ = strconv.Atoi(args[3])
	}

	decryptEntry := func(idx int) (crypto.EnvFile, error) {
		entry, err := crypto.EntryByVersion(vault, idx)
		if err != nil {
			return nil, err
		}
		plain, err := crypto.Decrypt(privKey, entry.Ciphertext)
		if err != nil {
			return nil, err
		}
		return crypto.ParseEnvBytes(plain)
	}

	oldEnv, err := decryptEntry(vA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decrypt v%d: %v\n", vA, err)
		os.Exit(1)
	}
	newEnv, err := decryptEntry(vB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decrypt v%d: %v\n", vB, err)
		os.Exit(1)
	}

	diff := crypto.DiffEnvFiles(oldEnv, newEnv)
	out := crypto.FormatDiff(diff)
	if out == "" {
		fmt.Println("no differences")
	} else {
		fmt.Print(out)
	}
}
