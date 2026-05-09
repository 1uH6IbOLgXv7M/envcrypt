package crypto

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadAuditLog_Missing(t *testing.T) {
	dir := t.TempDir()
	log, err := LoadAuditLog(dir)
	if err != nil {
		t.Fatalf("expected no error for missing log, got: %v", err)
	}
	if len(log.Events) != 0 {
		t.Fatalf("expected empty events, got %d", len(log.Events))
	}
}

func TestSaveAndLoadAuditLog(t *testing.T) {
	dir := t.TempDir()
	log := &AuditLog{
		Events: []AuditEvent{
			{Timestamp: time.Now().UTC(), Action: "push", Actor: "alice", Details: "v1"},
		},
	}
	if err := SaveAuditLog(dir, log); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadAuditLog(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(loaded.Events))
	}
	if loaded.Events[0].Action != "push" {
		t.Errorf("expected action 'push', got %q", loaded.Events[0].Action)
	}
	if loaded.Events[0].Actor != "alice" {
		t.Errorf("expected actor 'alice', got %q", loaded.Events[0].Actor)
	}
}

func TestAppendEvent(t *testing.T) {
	dir := t.TempDir()
	if err := AppendEvent(dir, "push", "bob", "v1"); err != nil {
		t.Fatalf("append 1: %v", err)
	}
	if err := AppendEvent(dir, "pull", "carol", "v1"); err != nil {
		t.Fatalf("append 2: %v", err)
	}
	log, err := LoadAuditLog(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(log.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(log.Events))
	}
	if log.Events[0].Action != "push" || log.Events[1].Action != "pull" {
		t.Errorf("unexpected actions: %v, %v", log.Events[0].Action, log.Events[1].Action)
	}
}

func TestAuditLogPath(t *testing.T) {
	path := AuditLogPath("/some/dir")
	expected := filepath.Join("/some/dir", ".envcrypt_audit.json")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

func TestSaveAuditLog_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	log := &AuditLog{}
	if err := SaveAuditLog(dir, log); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(AuditLogPath(dir))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
