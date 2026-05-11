package crypto

import (
	"os"
	"testing"
	"time"
)

func TestLoadSnapshotLog_Missing(t *testing.T) {
	dir := t.TempDir()
	log, err := LoadSnapshotLog(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(log.Snapshots) != 0 {
		t.Errorf("expected empty log, got %d entries", len(log.Snapshots))
	}
}

func TestSaveAndLoadSnapshotLog(t *testing.T) {
	dir := t.TempDir()
	snap := Snapshot{
		Version:   1,
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Keys:      []string{"DB_URL", "API_KEY"},
		Checksum:  "abc123",
	}
	log := &SnapshotLog{Snapshots: []Snapshot{snap}}
	if err := SaveSnapshotLog(dir, log); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadSnapshotLog(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Snapshots) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(loaded.Snapshots))
	}
	if loaded.Snapshots[0].Checksum != "abc123" {
		t.Errorf("checksum mismatch: %s", loaded.Snapshots[0].Checksum)
	}
}

func TestAppendSnapshot(t *testing.T) {
	dir := t.TempDir()
	for i := 1; i <= 3; i++ {
		snap := Snapshot{
			Version:   i,
			Timestamp: time.Now().UTC(),
			Keys:      []string{"KEY"},
			Checksum:  "chk",
		}
		if err := AppendSnapshot(dir, snap); err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
	}
	log, err := LoadSnapshotLog(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(log.Snapshots) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(log.Snapshots))
	}
}

func TestLatestSnapshot(t *testing.T) {
	log := &SnapshotLog{}
	if LatestSnapshot(log) != nil {
		t.Error("expected nil for empty log")
	}
	log.Snapshots = []Snapshot{
		{Version: 1, Checksum: "first"},
		{Version: 2, Checksum: "second"},
	}
	latest := LatestSnapshot(log)
	if latest == nil || latest.Checksum != "second" {
		t.Errorf("expected 'second', got %v", latest)
	}
}

func TestSnapshotLogPath(t *testing.T) {
	path := SnapshotLogPath("/some/dir")
	expected := "/some/dir/.envcrypt_snapshots.json"
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestLoadSnapshotLog_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := SnapshotLogPath(dir)
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := LoadSnapshotLog(dir)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
