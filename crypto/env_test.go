package crypto

import (
	"os"
	"testing"
)

const sampleEnv = `# Database config
DB_HOST=localhost
DB_PORT=5432
DB_NAME="myapp"

# App config
APP_SECRET='supersecret'
DEBUG=true
`

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseEnvFile_Basic(t *testing.T) {
	path := writeTempEnv(t, sampleEnv)
	ef, err := ParseEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", m["DB_HOST"])
	}
	if m["DB_NAME"] != "myapp" {
		t.Errorf("expected DB_NAME=myapp (unquoted), got %q", m["DB_NAME"])
	}
	if m["APP_SECRET"] != "supersecret" {
		t.Errorf("expected APP_SECRET=supersecret, got %q", m["APP_SECRET"])
	}
}

func TestParseEnvFile_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, sampleEnv)
	ef, err := ParseEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	comments, blanks := 0, 0
	for _, e := range ef.Entries {
		if e.Comment {
			comments++
		}
		if e.Blank {
			blanks++
		}
	}
	if comments != 2 {
		t.Errorf("expected 2 comment lines, got %d", comments)
	}
	if blanks < 1 {
		t.Errorf("expected at least 1 blank line, got %d", blanks)
	}
}

func TestParseEnvFile_MissingFile(t *testing.T) {
	_, err := ParseEnvFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestEnvFile_Serialize(t *testing.T) {
	path := writeTempEnv(t, sampleEnv)
	ef, err := ParseEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	serialized := ef.Serialize()
	if len(serialized) == 0 {
		t.Error("expected non-empty serialized output")
	}
	// Round-trip: re-parse the serialized content
	path2 := writeTempEnv(t, serialized)
	ef2, err := ParseEnvFile(path2)
	if err != nil {
		t.Fatalf("round-trip parse error: %v", err)
	}
	if len(ef2.Entries) != len(ef.Entries) {
		t.Errorf("entry count mismatch: want %d, got %d", len(ef.Entries), len(ef2.Entries))
	}
}
