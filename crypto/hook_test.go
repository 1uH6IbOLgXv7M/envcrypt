package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadHookStore_Missing(t *testing.T) {
	dir := t.TempDir()
	hs, err := LoadHookStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hs.Hooks) != 0 {
		t.Errorf("expected empty store, got %d hooks", len(hs.Hooks))
	}
}

func TestSaveAndLoadHookStore(t *testing.T) {
	dir := t.TempDir()
	hs := &HookStore{}
	hs.AddHook(HookPrePush, "make lint")
	hs.AddHook(HookPostPull, "echo pulled")

	if err := SaveHookStore(dir, hs); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadHookStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Hooks) != 2 {
		t.Errorf("expected 2 hooks, got %d", len(loaded.Hooks))
	}
}

func TestHookStore_AddAndGet(t *testing.T) {
	hs := &HookStore{}
	hs.AddHook(HookPrePush, "go test ./...")

	h := hs.GetHook(HookPrePush)
	if h == nil {
		t.Fatal("expected hook, got nil")
	}
	if h.Command != "go test ./..." {
		t.Errorf("unexpected command: %s", h.Command)
	}
}

func TestHookStore_AddHook_Overwrite(t *testing.T) {
	hs := &HookStore{}
	hs.AddHook(HookPrePush, "old command")
	hs.AddHook(HookPrePush, "new command")

	if len(hs.Hooks) != 1 {
		t.Errorf("expected 1 hook after overwrite, got %d", len(hs.Hooks))
	}
	if hs.Hooks[0].Command != "new command" {
		t.Errorf("expected updated command, got %s", hs.Hooks[0].Command)
	}
}

func TestHookStore_RemoveHook(t *testing.T) {
	hs := &HookStore{}
	hs.AddHook(HookPostPush, "echo done")

	if !hs.RemoveHook(HookPostPush) {
		t.Error("expected removal to return true")
	}
	if len(hs.Hooks) != 0 {
		t.Errorf("expected 0 hooks after removal, got %d", len(hs.Hooks))
	}
	if hs.RemoveHook(HookPostPush) {
		t.Error("expected removal of missing hook to return false")
	}
}

func TestHookStorePath(t *testing.T) {
	dir := "/tmp/project"
	expected := filepath.Join(dir, ".envcrypt", "hooks.json")
	if got := HookStorePath(dir); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestSaveHookStore_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	hs := &HookStore{}
	hs.AddHook(HookPrePull, "echo pre-pull")
	if err := SaveHookStore(dir, hs); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(HookStorePath(dir))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
