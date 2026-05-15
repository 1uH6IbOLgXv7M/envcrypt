package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChecksumEnvMap_Deterministic(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	c1 := ChecksumEnvMap(env)
	c2 := ChecksumEnvMap(env)
	if c1.Hash != c2.Hash {
		t.Errorf("expected deterministic hash, got %q vs %q", c1.Hash, c2.Hash)
	}
}

func TestChecksumEnvMap_OrderIndependent(t *testing.T) {
	a := map[string]string{"ALPHA": "1", "BETA": "2"}
	b := map[string]string{"BETA": "2", "ALPHA": "1"}
	if ChecksumEnvMap(a).Hash != ChecksumEnvMap(b).Hash {
		t.Error("checksum should be order-independent")
	}
}

func TestChecksumEnvMap_DifferentValues(t *testing.T) {
	a := map[string]string{"KEY": "value1"}
	b := map[string]string{"KEY": "value2"}
	if ChecksumEnvMap(a).Hash == ChecksumEnvMap(b).Hash {
		t.Error("different values should produce different hashes")
	}
}

func TestChecksumEnvMap_KeysPopulated(t *testing.T) {
	env := map[string]string{"Z": "1", "A": "2"}
	c := ChecksumEnvMap(env)
	if len(c.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(c.Keys))
	}
	if c.Keys[0] != "A" || c.Keys[1] != "Z" {
		t.Errorf("keys not sorted: %v", c.Keys)
	}
}

func TestChecksumEnvFile_Basic(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	_ = os.WriteFile(p, []byte("FOO=bar\nBAZ=qux\n"), 0600)

	c, err := ChecksumEnvFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Hash == "" {
		t.Error("expected non-empty hash")
	}
}

func TestChecksumEnvFile_Missing(t *testing.T) {
	_, err := ChecksumEnvFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestMatchChecksum_True(t *testing.T) {
	env := map[string]string{"X": "y"}
	c := ChecksumEnvMap(env)
	if !MatchChecksum(env, c.Hash) {
		t.Error("expected MatchChecksum to return true")
	}
}

func TestMatchChecksum_False(t *testing.T) {
	env := map[string]string{"X": "y"}
	if MatchChecksum(env, "deadbeef") {
		t.Error("expected MatchChecksum to return false for wrong hash")
	}
}
