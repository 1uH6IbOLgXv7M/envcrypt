package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupLintTestDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writeLintFile(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunLint_NoArgs(t *testing.T) {
	exited := false
	origExit := osExit
	osExit = func(code int) { exited = true; panic("exit") }
	defer func() {
		osExit = origExit
		recover()
	}()

	runLint([]string{})
	if !exited {
		t.Error("expected exit on no args")
	}
}

func TestRunLint_CleanFile(t *testing.T) {
	dir := setupLintTestDir(t)
	p := writeLintFile(t, dir, "DB_HOST=localhost\nDB_PORT=5432\n")

	// Should not panic or exit
	runLint([]string{p})
}

func TestRunLint_WithWarnings(t *testing.T) {
	dir := setupLintTestDir(t)
	p := writeLintFile(t, dir, "API_KEY=\nDB_HOST=localhost\n")

	// warnings only — should not call os.Exit
	runLint([]string{p})
}

func TestRunLint_MissingFile(t *testing.T) {
	exited := false
	origExit := osExit
	osExit = func(code int) { exited = true; panic("exit") }
	defer func() {
		osExit = origExit
		recover()
	}()

	runLint([]string{"/nonexistent/.env"})
	if !exited {
		t.Error("expected exit for missing file")
	}
}

func TestRunLint_WithErrors(t *testing.T) {
	dir := setupLintTestDir(t)
	p := writeLintFile(t, dir, "DB_HOST=a\nDB_HOST=b\n")

	exitCode := 0
	origExit := osExit
	osExit = func(code int) { exitCode = code; panic("exit") }
	defer func() {
		osExit = origExit
		recover()
	}()

	runLint([]string{p})
	if exitCode != 2 {
		t.Errorf("expected exit code 2 for errors, got %d", exitCode)
	}
}
