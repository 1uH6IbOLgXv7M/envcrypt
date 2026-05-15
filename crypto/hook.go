package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HookEvent represents a lifecycle event that triggers a hook.
type HookEvent string

const (
	HookPrePush  HookEvent = "pre-push"
	HookPostPush HookEvent = "post-push"
	HookPrePull  HookEvent = "pre-pull"
	HookPostPull HookEvent = "post-pull"
)

// HookConfig defines a hook script or command for a given event.
type HookConfig struct {
	Event   HookEvent `json:"event"`
	Command string    `json:"command"`
	Enabled bool      `json:"enabled"`
}

// HookStore holds all registered hooks.
type HookStore struct {
	Hooks     []HookConfig `json:"hooks"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// HookStorePath returns the path to the hook store file.
func HookStorePath(dir string) string {
	return filepath.Join(dir, ".envcrypt", "hooks.json")
}

// LoadHookStore loads the hook store from disk, returning an empty store if missing.
func LoadHookStore(dir string) (*HookStore, error) {
	path := HookStorePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &HookStore{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read hook store: %w", err)
	}
	var hs HookStore
	if err := json.Unmarshal(data, &hs); err != nil {
		return nil, fmt.Errorf("parse hook store: %w", err)
	}
	return &hs, nil
}

// SaveHookStore persists the hook store to disk.
func SaveHookStore(dir string, hs *HookStore) error {
	path := HookStorePath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("mkdir hook store: %w", err)
	}
	hs.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(hs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal hook store: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// AddHook registers or replaces a hook for the given event.
func (hs *HookStore) AddHook(event HookEvent, command string) {
	for i, h := range hs.Hooks {
		if h.Event == event {
			hs.Hooks[i] = HookConfig{Event: event, Command: command, Enabled: true}
			return
		}
	}
	hs.Hooks = append(hs.Hooks, HookConfig{Event: event, Command: command, Enabled: true})
}

// RemoveHook removes the hook for the given event, returning true if found.
func (hs *HookStore) RemoveHook(event HookEvent) bool {
	for i, h := range hs.Hooks {
		if h.Event == event {
			hs.Hooks = append(hs.Hooks[:i], hs.Hooks[i+1:]...)
			return true
		}
	}
	return false
}

// GetHook returns the hook config for the given event, or nil if not found.
func (hs *HookStore) GetHook(event HookEvent) *HookConfig {
	for i, h := range hs.Hooks {
		if h.Event == event && h.Enabled {
			return &hs.Hooks[i]
		}
	}
	return nil
}
