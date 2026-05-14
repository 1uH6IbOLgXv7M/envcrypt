package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTagStore_Missing(t *testing.T) {
	dir := t.TempDir()
	store, err := LoadTagStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Tags) != 0 {
		t.Errorf("expected empty store, got %d tags", len(store.Tags))
	}
}

func TestSaveAndLoadTagStore(t *testing.T) {
	dir := t.TempDir()
	store := &TagStore{}
	store.AddTag("v1.0", 3, "release")
	if err := SaveTagStore(dir, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadTagStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(loaded.Tags))
	}
	if loaded.Tags[0].Name != "v1.0" || loaded.Tags[0].Version != 3 {
		t.Errorf("tag mismatch: %+v", loaded.Tags[0])
	}
}

func TestTagStore_AddAndGet(t *testing.T) {
	store := &TagStore{}
	store.AddTag("stable", 5, "stable release")
	tag := store.GetTag("stable")
	if tag == nil {
		t.Fatal("expected tag, got nil")
	}
	if tag.Version != 5 || tag.Message != "stable release" {
		t.Errorf("unexpected tag: %+v", tag)
	}
}

func TestTagStore_AddTag_Overwrite(t *testing.T) {
	store := &TagStore{}
	store.AddTag("latest", 1, "")
	store.AddTag("latest", 7, "updated")
	if len(store.Tags) != 1 {
		t.Errorf("expected 1 tag after overwrite, got %d", len(store.Tags))
	}
	if store.Tags[0].Version != 7 {
		t.Errorf("expected version 7, got %d", store.Tags[0].Version)
	}
}

func TestTagStore_RemoveTag(t *testing.T) {
	store := &TagStore{}
	store.AddTag("old", 2, "")
	if !store.RemoveTag("old") {
		t.Error("expected RemoveTag to return true")
	}
	if store.GetTag("old") != nil {
		t.Error("expected tag to be removed")
	}
	if store.RemoveTag("old") {
		t.Error("expected RemoveTag to return false for missing tag")
	}
}

func TestTagStorePath(t *testing.T) {
	path := TagStorePath("/some/dir")
	expected := filepath.Join("/some/dir", ".envcrypt_tags.json")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

func TestLoadTagStore_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, ".envcrypt_tags.json"), []byte("not-json"), 0644)
	_, err := LoadTagStore(dir)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
