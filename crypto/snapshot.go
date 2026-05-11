package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of an env file's metadata.
type Snapshot struct {
	Version   int               `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Keys      []string          `json:"keys"`
	Checksum  string            `json:"checksum"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// SnapshotLog holds all snapshots for a vault.
type SnapshotLog struct {
	Snapshots []Snapshot `json:"snapshots"`
}

// SnapshotLogPath returns the path to the snapshot log file.
func SnapshotLogPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_snapshots.json")
}

// LoadSnapshotLog loads the snapshot log from disk, returning empty if missing.
func LoadSnapshotLog(dir string) (*SnapshotLog, error) {
	path := SnapshotLogPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &SnapshotLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var log SnapshotLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &log, nil
}

// SaveSnapshotLog writes the snapshot log to disk.
func SaveSnapshotLog(dir string, log *SnapshotLog) error {
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := SnapshotLogPath(dir)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// AppendSnapshot adds a new snapshot entry to the log.
func AppendSnapshot(dir string, snap Snapshot) error {
	log, err := LoadSnapshotLog(dir)
	if err != nil {
		return err
	}
	log.Snapshots = append(log.Snapshots, snap)
	return SaveSnapshotLog(dir, log)
}

// LatestSnapshot returns the most recent snapshot, or nil if none exist.
func LatestSnapshot(log *SnapshotLog) *Snapshot {
	if len(log.Snapshots) == 0 {
		return nil
	}
	s := log.Snapshots[len(log.Snapshots)-1]
	return &s
}
