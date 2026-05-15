package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"envcrypt/crypto"
)

func setupHookTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envcrypt"), 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	return dir
}

func TestRunHookAdd_NoArgs(t *testing.T) {
	err := runHookAdd([]string{})
	if err == nil {
		t.Error("expected error for no args")
	}
}

func TestRunHookAdd_And_List(t *testing.T) {
	dir := setupHookTestDir(t)
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	if err := runHookAdd([]string{"pre-push", "make", "lint"}); err != nil {
		t.Fatalf("add: %v", err)
	}

	hs, err := crypto.LoadHookStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(hs.Hooks) != 1 {
		t.Errorf("expected 1 hook, got %d", len(hs.Hooks))
	}
	if hs.Hooks[0].Command != "make lint" {
		t.Errorf("unexpected command: %s", hs.Hooks[0].Command)
	}

	if err := runHookList([]string{}); err != nil {
		t.Fatalf("list: %v", err)
	}
}

func TestRunHookRemove(t *testing.T) {
	dir := setupHookTestDir(t)
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	_ = runHookAdd([]string{"post-push", "echo done"})

	if err := runHookRemove([]string{"post-push"}); err != nil {
		t.Fatalf("remove: %v", err)
	}

	hs, _ := crypto.LoadHookStore(dir)
	if len(hs.Hooks) != 0 {
		t.Errorf("expected 0 hooks after removal, got %d", len(hs.Hooks))
	}
}

func TestRunHookRemove_NoArgs(t *testing.T) {
	err := runHookRemove([]string{})
	if err == nil {
		t.Error("expected error for no args")
	}
}

func TestRunHookRemove_NotFound(t *testing.T) {
	dir := setupHookTestDir(t)
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	err := runHookRemove([]string{"pre-pull"})
	if err == nil {
		t.Error("expected error when hook not found")
	}
}

func TestRunHookList_Empty(t *testing.T) {
	dir := setupHookTestDir(t)
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	if err := runHookList([]string{}); err != nil {
		t.Fatalf("list empty: %v", err)
	}
}
