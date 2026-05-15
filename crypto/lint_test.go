package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func writeLintEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLintEnvFile_Clean(t *testing.T) {
	p := writeLintEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	results, err := LintEnvFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no findings, got %d", len(results))
	}
}

func TestLintEnvFile_EmptyValue(t *testing.T) {
	p := writeLintEnv(t, "API_KEY=\nDB_HOST=localhost\n")
	results, err := LintEnvFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(results))
	}
	if results[0].Severity != "warn" || results[0].Key != "API_KEY" {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

func TestLintEnvFile_DuplicateKey(t *testing.T) {
	p := writeLintEnv(t, "DB_HOST=localhost\nDB_HOST=remotehost\n")
	results, err := LintEnvFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(results))
	}
	if results[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", results[0].Severity)
	}
}

func TestLintEnvFile_BadKeyCase(t *testing.T) {
	p := writeLintEnv(t, "dbHost=localhost\n")
	results, err := LintEnvFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(results))
	}
	if results[0].Message != "key should be UPPER_SNAKE_CASE" {
		t.Errorf("unexpected message: %s", results[0].Message)
	}
}

func TestLintEnvFile_MissingFile(t *testing.T) {
	_, err := LintEnvFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestIsUpperSnakeCase(t *testing.T) {
	cases := []struct {
		input string
		ok    bool
	}{
		{"DB_HOST", true},
		{"API_KEY_123", true},
		{"dbHost", false},
		{"DB-HOST", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isUpperSnakeCase(c.input); got != c.ok {
			t.Errorf("isUpperSnakeCase(%q) = %v, want %v", c.input, got, c.ok)
		}
	}
}
