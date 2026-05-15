package crypto

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempWatch(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempWatch: %v", err)
	}
	return p
}

func TestWatchEnvFile_NoChangeNoEvent(t *testing.T) {
	path := writeTempWatch(t, "KEY=value\n")
	done := make(chan struct{})
	ch, err := WatchEnvFile(path, 20*time.Millisecond, done)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(60 * time.Millisecond)
	close(done)
	select {
	case ev, ok := <-ch:
		if ok {
			t.Fatalf("expected no event, got %+v", ev)
		}
	default:
	}
}

func TestWatchEnvFile_DetectsChange(t *testing.T) {
	path := writeTempWatch(t, "KEY=original\n")
	done := make(chan struct{})
	ch, err := WatchEnvFile(path, 20*time.Millisecond, done)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer close(done)

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("path mismatch: got %s", ev.Path)
		}
		if ev.Checksum == "" {
			t.Error("expected non-empty checksum")
		}
		if ev.DetectedAt.IsZero() {
			t.Error("expected non-zero DetectedAt")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatchEnvFile_MissingFile(t *testing.T) {
	_, err := WatchEnvFile("/nonexistent/.env", 20*time.Millisecond, make(chan struct{}))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFileChecksum_Deterministic(t *testing.T) {
	path := writeTempWatch(t, "A=1\nB=2\n")
	c1, err := fileChecksum(path)
	if err != nil {
		t.Fatalf("checksum: %v", err)
	}
	c2, err := fileChecksum(path)
	if err != nil {
		t.Fatalf("checksum: %v", err)
	}
	if c1 != c2 {
		t.Errorf("expected deterministic checksum, got %s vs %s", c1, c2)
	}
}
