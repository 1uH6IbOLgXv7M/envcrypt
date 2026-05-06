package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// VaultEntry represents a single encrypted .env snapshot.
type VaultEntry struct {
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
	Ciphertext []byte `json:"ciphertext"`
}

// Vault holds a versioned list of encrypted env snapshots.
type Vault struct {
	Entries []VaultEntry `json:"entries"`
}

// LoadVault reads a vault from disk, or returns an empty one if not found.
func LoadVault(path string) (*Vault, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Vault{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read vault: %w", err)
	}
	var v Vault
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parse vault: %w", err)
	}
	return &v, nil
}

// SaveVault writes the vault to disk as JSON.
func SaveVault(path string, v *Vault) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal vault: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write vault: %w", err)
	}
	return nil
}

// AddEntry appends a new encrypted snapshot to the vault.
func (v *Vault) AddEntry(ciphertext []byte) VaultEntry {
	nextVersion := 1
	if len(v.Entries) > 0 {
		nextVersion = v.Entries[len(v.Entries)-1].Version + 1
	}
	entry := VaultEntry{
		Version:    nextVersion,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Ciphertext: ciphertext,
	}
	v.Entries = append(v.Entries, entry)
	return entry
}

// LatestEntry returns the most recent vault entry, or an error if empty.
func (v *Vault) LatestEntry() (VaultEntry, error) {
	if len(v.Entries) == 0 {
		return VaultEntry{}, fmt.Errorf("vault is empty")
	}
	return v.Entries[len(v.Entries)-1], nil
}

// EntryByVersion finds a vault entry by version number.
func (v *Vault) EntryByVersion(version int) (VaultEntry, error) {
	for _, e := range v.Entries {
		if e.Version == version {
			return e, nil
		}
	}
	return VaultEntry{}, fmt.Errorf("version %d not found in vault", version)
}
