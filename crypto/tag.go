package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Tag represents a named pointer to a specific vault version.
type Tag struct {
	Name      string    `json:"name"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message,omitempty"`
}

// TagStore holds all tags for a vault.
type TagStore struct {
	Tags []Tag `json:"tags"`
}

// TagStorePath returns the path to the tag store file.
func TagStorePath(dir string) string {
	return filepath.Join(dir, ".envcrypt_tags.json")
}

// LoadTagStore loads the tag store from disk, returning an empty store if missing.
func LoadTagStore(dir string) (*TagStore, error) {
	path := TagStorePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &TagStore{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read tag store: %w", err)
	}
	var store TagStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse tag store: %w", err)
	}
	return &store, nil
}

// SaveTagStore writes the tag store to disk.
func SaveTagStore(dir string, store *TagStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tag store: %w", err)
	}
	return os.WriteFile(TagStorePath(dir), data, 0644)
}

// AddTag adds or overwrites a tag in the store.
func (s *TagStore) AddTag(name string, version int, message string) {
	for i, t := range s.Tags {
		if t.Name == name {
			s.Tags[i] = Tag{Name: name, Version: version, CreatedAt: time.Now(), Message: message}
			return
		}
	}
	s.Tags = append(s.Tags, Tag{Name: name, Version: version, CreatedAt: time.Now(), Message: message})
}

// GetTag returns a tag by name, or nil if not found.
func (s *TagStore) GetTag(name string) *Tag {
	for _, t := range s.Tags {
		if t.Name == name {
			copy := t
			return &copy
		}
	}
	return nil
}

// RemoveTag removes a tag by name. Returns false if not found.
func (s *TagStore) RemoveTag(name string) bool {
	for i, t := range s.Tags {
		if t.Name == name {
			s.Tags = append(s.Tags[:i], s.Tags[i+1:]...)
			return true
		}
	}
	return false
}
