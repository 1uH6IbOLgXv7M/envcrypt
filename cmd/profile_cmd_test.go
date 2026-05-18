package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"envcrypt/crypto"
)

func setupProfileTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envcrypt"), 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	return dir
}

func TestRunProfileAdd_NoArgs(t *testing.T) {
	if err := runProfileAdd([]string{}); err == nil {
		t.Error("expected error for no args")
	}
	if err := runProfileAdd([]string{"dev"}); err == nil {
		t.Error("expected error for one arg")
	}
}

func TestRunProfileAdd_And_List(t *testing.T) {
	dir := setupProfileTestDir(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	if err := runProfileAdd([]string{"dev", ".env.dev"}); err != nil {
		t.Fatalf("add: %v", err)
	}
	store, err := crypto.LoadProfileStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, ok := store.GetProfile("dev"); !ok {
		t.Error("expected profile 'dev' to exist")
	}
	if err := runProfileList([]string{}); err != nil {
		t.Fatalf("list: %v", err)
	}
}

func TestRunProfileRemove(t *testing.T) {
	dir := setupProfileTestDir(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	runProfileAdd([]string{"staging", ".env.staging"})
	if err := runProfileRemove([]string{"staging"}); err != nil {
		t.Fatalf("remove: %v", err)
	}
	store, _ := crypto.LoadProfileStore(dir)
	if _, ok := store.GetProfile("staging"); ok {
		t.Error("expected profile to be removed")
	}
}

func TestRunProfileRemove_NoArgs(t *testing.T) {
	if err := runProfileRemove([]string{}); err == nil {
		t.Error("expected error for no args")
	}
}

func TestRunProfileRemove_NotFound(t *testing.T) {
	dir := setupProfileTestDir(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	if err := runProfileRemove([]string{"nonexistent"}); err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestRunProfileList_Empty(t *testing.T) {
	dir := setupProfileTestDir(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	if err := runProfileList([]string{}); err != nil {
		t.Fatalf("list on empty store: %v", err)
	}
}
