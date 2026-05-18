package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Profile represents a named environment configuration profile.
type Profile struct {
	Name      string            `json:"name"`
	EnvFile   string            `json:"env_file"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ProfileStore holds all named profiles.
type ProfileStore struct {
	Profiles map[string]Profile `json:"profiles"`
}

// ProfileStorePath returns the path to the profile store file.
func ProfileStorePath(dir string) string {
	return filepath.Join(dir, ".envcrypt", "profiles.json")
}

// LoadProfileStore loads the profile store from disk, returning an empty store if missing.
func LoadProfileStore(dir string) (*ProfileStore, error) {
	path := ProfileStorePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &ProfileStore{Profiles: make(map[string]Profile)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read profile store: %w", err)
	}
	var store ProfileStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse profile store: %w", err)
	}
	if store.Profiles == nil {
		store.Profiles = make(map[string]Profile)
	}
	return &store, nil
}

// SaveProfileStore persists the profile store to disk.
func SaveProfileStore(dir string, store *ProfileStore) error {
	path := ProfileStorePath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("mkdir profile store: %w", err)
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profile store: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// AddProfile adds or updates a profile in the store.
func (s *ProfileStore) AddProfile(p Profile) {
	now := time.Now().UTC()
	if existing, ok := s.Profiles[p.Name]; ok {
		p.CreatedAt = existing.CreatedAt
	} else {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
	s.Profiles[p.Name] = p
}

// GetProfile retrieves a profile by name.
func (s *ProfileStore) GetProfile(name string) (Profile, bool) {
	p, ok := s.Profiles[name]
	return p, ok
}

// RemoveProfile deletes a profile by name.
func (s *ProfileStore) RemoveProfile(name string) bool {
	if _, ok := s.Profiles[name]; !ok {
		return false
	}
	delete(s.Profiles, name)
	return true
}
