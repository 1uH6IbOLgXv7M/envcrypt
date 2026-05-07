package crypto

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// KeyRingEntry represents a named key pair stored in the keyring.
type KeyRingEntry struct {
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
}

// KeyRing holds a collection of named public keys.
type KeyRing struct {
	Entries []KeyRingEntry `json:"entries"`
}

// KeyRingPath returns the default path for the keyring file.
func KeyRingPath(dir string) string {
	return filepath.Join(dir, "keyring.json")
}

// LoadKeyRing loads a keyring from disk, returning an empty one if missing.
func LoadKeyRing(path string) (*KeyRing, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &KeyRing{}, nil
	}
	if err != nil {
		return nil, err
	}
	var kr KeyRing
	if err := json.Unmarshal(data, &kr); err != nil {
		return nil, err
	}
	return &kr, nil
}

// SaveKeyRing persists the keyring to disk.
func SaveKeyRing(path string, kr *KeyRing) error {
	data, err := json.MarshalIndent(kr, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// AddKey adds or replaces a named key in the keyring.
func (kr *KeyRing) AddKey(name, publicKey string) {
	for i, e := range kr.Entries {
		if e.Name == name {
			kr.Entries[i].PublicKey = publicKey
			kr.Entries[i].CreatedAt = time.Now().UTC()
			return
		}
	}
	kr.Entries = append(kr.Entries, KeyRingEntry{
		Name:      name,
		PublicKey: publicKey,
		CreatedAt: time.Now().UTC(),
	})
}

// GetKey retrieves a public key by name.
func (kr *KeyRing) GetKey(name string) (string, bool) {
	for _, e := range kr.Entries {
		if e.Name == name {
			return e.PublicKey, true
		}
	}
	return "", false
}

// RemoveKey removes a named key from the keyring.
func (kr *KeyRing) RemoveKey(name string) bool {
	for i, e := range kr.Entries {
		if e.Name == name {
			kr.Entries = append(kr.Entries[:i], kr.Entries[i+1:]...)
			return true
		}
	}
	return false
}
