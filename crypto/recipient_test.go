package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRecipientsFile_Missing(t *testing.T) {
	rf, err := LoadRecipientsFile("/nonexistent/path/.env-recipients")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(rf.Recipients) != 0 {
		t.Errorf("expected empty recipients, got %d", len(rf.Recipients))
	}
}

func TestSaveAndLoadRecipientsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env-recipients")

	rf := &RecipientsFile{}
	rf.AddRecipient("alice", "age1alicepubkey")
	rf.AddRecipient("bob", "age1bobpubkey")

	if err := SaveRecipientsFile(path, rf); err != nil {
		t.Fatalf("SaveRecipientsFile: %v", err)
	}

	loaded, err := LoadRecipientsFile(path)
	if err != nil {
		t.Fatalf("LoadRecipientsFile: %v", err)
	}
	if len(loaded.Recipients) != 2 {
		t.Fatalf("expected 2 recipients, got %d", len(loaded.Recipients))
	}
	if loaded.Recipients[0].Alias != "alice" || loaded.Recipients[0].PublicKey != "age1alicepubkey" {
		t.Errorf("unexpected first recipient: %+v", loaded.Recipients[0])
	}
}

func TestAddRecipient_Overwrite(t *testing.T) {
	rf := &RecipientsFile{}
	rf.AddRecipient("alice", "age1old")
	rf.AddRecipient("alice", "age1new")
	if len(rf.Recipients) != 1 {
		t.Errorf("expected 1 recipient after overwrite, got %d", len(rf.Recipients))
	}
	if rf.Recipients[0].PublicKey != "age1new" {
		t.Errorf("expected updated key, got %q", rf.Recipients[0].PublicKey)
	}
}

func TestGetRecipient_NotFound(t *testing.T) {
	rf := &RecipientsFile{}
	_, err := rf.GetRecipient("nobody")
	if err == nil {
		t.Error("expected error for missing recipient")
	}
}

func TestRemoveRecipient(t *testing.T) {
	rf := &RecipientsFile{}
	rf.AddRecipient("alice", "age1alicepubkey")
	rf.AddRecipient("bob", "age1bobpubkey")

	removed := rf.RemoveRecipient("alice")
	if !removed {
		t.Error("expected true when removing existing recipient")
	}
	if len(rf.Recipients) != 1 {
		t.Errorf("expected 1 recipient after removal, got %d", len(rf.Recipients))
	}
	notRemoved := rf.RemoveRecipient("nobody")
	if notRemoved {
		t.Error("expected false when removing non-existent recipient")
	}
}

func TestLoadRecipientsFile_MalformedLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env-recipients")
	_ = os.WriteFile(path, []byte("badline\n"), 0644)
	_, err := LoadRecipientsFile(path)
	if err == nil {
		t.Error("expected error for malformed line")
	}
}
