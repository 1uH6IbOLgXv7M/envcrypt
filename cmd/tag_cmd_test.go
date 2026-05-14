package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"envcrypt/crypto"
)

func setupTagTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Generate keys and push a couple of vault entries.
	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	if err := crypto.SaveKeyPair(dir, pub, priv); err != nil {
		t.Fatalf("save keys: %v", err)
	}
	envPath := filepath.Join(dir, ".env")
	_ = os.WriteFile(envPath, []byte("KEY=value\n"), 0644)
	if err := runVaultPush([]string{envPath}, dir); err != nil {
		t.Fatalf("vault push: %v", err)
	}
	return dir
}

func TestRunTagAdd_NoArgs(t *testing.T) {
	dir := t.TempDir()
	if err := runTagAdd([]string{}, dir); err == nil {
		t.Error("expected error for no args")
	}
}

func TestRunTagAdd_InvalidVersion(t *testing.T) {
	dir := setupTagTestDir(t)
	if err := runTagAdd([]string{"v1", "notanumber"}, dir); err == nil {
		t.Error("expected error for invalid version")
	}
}

func TestRunTagAdd_VersionNotInVault(t *testing.T) {
	dir := setupTagTestDir(t)
	if err := runTagAdd([]string{"v99", "99"}, dir); err == nil {
		t.Error("expected error for missing vault version")
	}
}

func TestRunTagAdd_And_List(t *testing.T) {
	dir := setupTagTestDir(t)
	if err := runTagAdd([]string{"release", "1", "first release"}, dir); err != nil {
		t.Fatalf("tag add: %v", err)
	}
	store, err := crypto.LoadTagStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(store.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(store.Tags))
	}
	if store.Tags[0].Name != "release" {
		t.Errorf("unexpected tag name: %s", store.Tags[0].Name)
	}
	if err := runTagList(dir); err != nil {
		t.Errorf("tag list: %v", err)
	}
}

func TestRunTagList_Empty(t *testing.T) {
	dir := t.TempDir()
	if err := runTagList(dir); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTagRemove_NoArgs(t *testing.T) {
	dir := t.TempDir()
	if err := runTagRemove([]string{}, dir); err == nil {
		t.Error("expected error for no args")
	}
}

func TestRunTagRemove(t *testing.T) {
	dir := setupTagTestDir(t)
	_ = runTagAdd([]string{"beta", "1", ""}, dir)
	if err := runTagRemove([]string{"beta"}, dir); err != nil {
		t.Fatalf("tag remove: %v", err)
	}
	store, _ := crypto.LoadTagStore(dir)
	if store.GetTag("beta") != nil {
		t.Error("expected tag to be removed")
	}
}

func TestRunTagRemove_NotFound(t *testing.T) {
	dir := t.TempDir()
	if err := runTagRemove([]string{"ghost"}, dir); err == nil {
		t.Error("expected error for missing tag")
	}
}
