package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPolicy_Missing(t *testing.T) {
	dir := t.TempDir()
	p, err := LoadPolicy(dir)
	if err != nil {
		t.Fatalf("expected default policy, got error: %v", err)
	}
	if p.MinRecipients != 1 {
		t.Errorf("expected default MinRecipients=1, got %d", p.MinRecipients)
	}
	if p.MaxVersions != 100 {
		t.Errorf("expected default MaxVersions=100, got %d", p.MaxVersions)
	}
}

func TestSaveAndLoadPolicy(t *testing.T) {
	dir := t.TempDir()
	orig := &Policy{
		MinRecipients:   2,
		RequireSign:     true,
		AllowedKeys:     []string{"alice", "bob"},
		MaxVersions:     50,
		RequireAuditLog: true,
	}
	if err := SavePolicy(dir, orig); err != nil {
		t.Fatalf("save policy: %v", err)
	}
	loaded, err := LoadPolicy(dir)
	if err != nil {
		t.Fatalf("load policy: %v", err)
	}
	if loaded.MinRecipients != orig.MinRecipients {
		t.Errorf("MinRecipients mismatch: got %d", loaded.MinRecipients)
	}
	if !loaded.RequireSign {
		t.Error("expected RequireSign=true")
	}
	if len(loaded.AllowedKeys) != 2 {
		t.Errorf("expected 2 allowed keys, got %d", len(loaded.AllowedKeys))
	}
}

func TestPolicyPath(t *testing.T) {
	dir := "/tmp/myproject"
	expected := filepath.Join(dir, ".envcrypt", "policy.json")
	if got := PolicyPath(dir); got != expected {
		t.Errorf("PolicyPath: got %q, want %q", got, expected)
	}
}

func TestPolicy_Validate_MinRecipients(t *testing.T) {
	p := &Policy{MinRecipients: 3, MaxVersions: 100}
	v := &Vault{}
	kr := &KeyRing{Keys: map[string]string{"a": "k1", "b": "k2"}}
	if err := p.Validate(v, kr); err == nil {
		t.Error("expected error for insufficient recipients")
	}
}

func TestPolicy_Validate_AllowedKeys(t *testing.T) {
	p := &Policy{MinRecipients: 1, MaxVersions: 100, AllowedKeys: []string{"alice"}}
	v := &Vault{}
	kr := &KeyRing{Keys: map[string]string{"alice": "k1", "eve": "k2"}}
	if err := p.Validate(v, kr); err == nil {
		t.Error("expected error for disallowed key 'eve'")
	}
}

func TestPolicy_Validate_MaxVersions(t *testing.T) {
	p := &Policy{MinRecipients: 1, MaxVersions: 2}
	v := &Vault{Entries: []VaultEntry{{Version: 1}, {Version: 2}, {Version: 3}}}
	kr := &KeyRing{Keys: map[string]string{"a": "k"}}
	if err := p.Validate(v, kr); err == nil {
		t.Error("expected error for exceeding max versions")
	}
}

func TestPolicy_Validate_Pass(t *testing.T) {
	p := &Policy{MinRecipients: 1, MaxVersions: 10, AllowedKeys: []string{"alice", "bob"}}
	v := &Vault{Entries: []VaultEntry{{Version: 1}}}
	kr := &KeyRing{Keys: map[string]string{"alice": "k1", "bob": "k2"}}
	if err := p.Validate(v, kr); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestSavePolicy_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	p := defaultPolicy()
	if err := SavePolicy(dir, p); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(PolicyPath(dir))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
