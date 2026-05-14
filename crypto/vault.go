package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// VaultEntry holds one encrypted version of an env file.
type VaultEntry struct {
	Version   int       `json:"version"`
	Ciphertext []byte   `json:"ciphertext"`
	CreatedAt time.Time `json:"created_at"`
}

// Vault is the on-disk encrypted store.
type Vault struct {
	Entries []VaultEntry `json:"entries"`
}

func vaultPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_vault.json")
}

// LoadVault reads the vault from disk, returning an empty vault if missing.
func LoadVault(dir string) (*Vault, error) {
	data, err := os.ReadFile(vaultPath(dir))
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

// SaveVault writes the vault to disk.
func SaveVault(dir string, v *Vault) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal vault: %w", err)
	}
	return os.WriteFile(vaultPath(dir), data, 0644)
}

// AddEntry appends a new encrypted entry, auto-incrementing the version.
func (v *Vault) AddEntry(ciphertext []byte) VaultEntry {
	next := 1
	if len(v.Entries) > 0 {
		next = v.Entries[len(v.Entries)-1].Version + 1
	}
	entry := VaultEntry{Version: next, Ciphertext: ciphertext, CreatedAt: time.Now()}
	v.Entries = append(v.Entries, entry)
	return entry
}

// LatestEntry returns the most recently added entry.
func (v *Vault) LatestEntry() (*VaultEntry, error) {
	if len(v.Entries) == 0 {
		return nil, fmt.Errorf("vault is empty")
	}
	e := v.Entries[len(v.Entries)-1]
	return &e, nil
}

// EntryByVersion returns the entry matching the given version number.
func (v *Vault) EntryByVersion(version int) (*VaultEntry, error) {
	for _, e := range v.Entries {
		if e.Version == version {
			copy := e
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("version %d not found", version)
}
