package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupSchemaTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func writeSchemaEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func writeSchemaJSON(t *testing.T, dir string, fields []map[string]interface{}) string {
	t.Helper()
	data, _ := json.Marshal(map[string]interface{}{"fields": fields})
	p := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(p, data, 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunSchemaValidate_NoArgs(t *testing.T) {
	err := runSchemaValidate([]string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRunSchemaValidate_Pass(t *testing.T) {
	dir := setupSchemaTestDir(t)
	envFile := writeSchemaEnv(t, dir, ".env", "HOST=localhost\nPORT=8080\n")
	schemaFile := writeSchemaJSON(t, dir, []map[string]interface{}{
		{"key": "HOST", "required": true},
		{"key": "PORT", "required": true, "pattern": `^\d+$`},
	})
	if err := runSchemaValidate([]string{envFile, schemaFile}); err != nil {
		t.Fatalf("expected pass, got error: %v", err)
	}
}

func TestRunSchemaValidate_Fail_MissingKey(t *testing.T) {
	dir := setupSchemaTestDir(t)
	envFile := writeSchemaEnv(t, dir, ".env", "HOST=localhost\n")
	schemaFile := writeSchemaJSON(t, dir, []map[string]interface{}{
		{"key": "HOST", "required": true},
		{"key": "PORT", "required": true},
	})
	if err := runSchemaValidate([]string{envFile, schemaFile}); err == nil {
		t.Fatal("expected validation error for missing PORT")
	}
}

func TestRunSchemaValidate_Fail_PatternMismatch(t *testing.T) {
	dir := setupSchemaTestDir(t)
	envFile := writeSchemaEnv(t, dir, ".env", "PORT=notanumber\n")
	schemaFile := writeSchemaJSON(t, dir, []map[string]interface{}{
		{"key": "PORT", "required": true, "pattern": `^\d+$`},
	})
	if err := runSchemaValidate([]string{envFile, schemaFile}); err == nil {
		t.Fatal("expected validation error for pattern mismatch")
	}
}

func TestRunSchemaGenerate_NoArgs(t *testing.T) {
	if err := runSchemaGenerate([]string{}); err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRunSchemaGenerate_Output(t *testing.T) {
	dir := setupSchemaTestDir(t)
	envFile := writeSchemaEnv(t, dir, ".env", "DB_HOST=localhost\nDB_PORT=5432\n")
	if err := runSchemaGenerate([]string{envFile}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadSchemaFile_Missing(t *testing.T) {
	_, err := loadSchemaFile("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing schema file")
	}
}
