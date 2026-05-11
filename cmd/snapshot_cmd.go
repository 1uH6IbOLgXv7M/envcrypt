package cmd

import (
	"fmt"
	"os"
	"time"

	"envcrypt/crypto"
)

// runSnapshotList lists all snapshots recorded in the current directory.
func runSnapshotList(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	log, err := crypto.LoadSnapshotLog(dir)
	if err != nil {
		return fmt.Errorf("snapshot list: %w", err)
	}
	if len(log.Snapshots) == 0 {
		fmt.Println("No snapshots found.")
		return nil
	}
	fmt.Printf("%-8s %-30s %-16s %s\n", "VERSION", "TIMESTAMP", "CHECKSUM", "KEYS")
	for _, s := range log.Snapshots {
		keys := fmt.Sprintf("%d key(s)", len(s.Keys))
		chk := s.Checksum
		if len(chk) > 12 {
			chk = chk[:12] + "..."
		}
		fmt.Printf("%-8d %-30s %-16s %s\n",
			s.Version,
			s.Timestamp.Format(time.RFC3339),
			chk,
			keys,
		)
	}
	return nil
}

// runSnapshotCreate captures a snapshot of the current vault state.
func runSnapshotCreate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: snapshot create <dir> <private-key-file>")
	}
	dir := args[0]
	privKeyFile := args[1]

	privKey, err := crypto.LoadPrivateKey(privKeyFile)
	if err != nil {
		return fmt.Errorf("snapshot create: load key: %w", err)
	}

	vault, err := crypto.LoadVault(dir)
	if err != nil {
		return fmt.Errorf("snapshot create: load vault: %w", err)
	}

	entry := crypto.LatestEntry(vault)
	if entry == nil {
		return fmt.Errorf("snapshot create: no vault entries found")
	}

	plain, err := crypto.Decrypt(privKey, entry.Ciphertext)
	if err != nil {
		return fmt.Errorf("snapshot create: decrypt: %w", err)
	}

	env, err := crypto.ParseEnvBytes(plain)
	if err != nil {
		return fmt.Errorf("snapshot create: parse env: %w", err)
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	snap := crypto.Snapshot{
		Version:   entry.Version,
		Timestamp: time.Now().UTC(),
		Keys:      keys,
		Checksum:  fmt.Sprintf("%x", len(plain)),
	}

	if err := crypto.AppendSnapshot(dir, snap); err != nil {
		return fmt.Errorf("snapshot create: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot created for version %d (%d keys).\n", snap.Version, len(keys))
	return nil
}
