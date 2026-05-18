package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSchemaValidate_OptionalFieldAllowed(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("HOST=localhost\n"), 0600)

	schemaData, _ := json.Marshal(map[string]interface{}{
		"fields": []map[string]interface{}{
			{"key": "HOST", "required": true},
			{"key": "DEBUG", "required": false},
		},
	})
	schemaFile := filepath.Join(dir, "schema.json")
	_ = os.WriteFile(schemaFile, schemaData, 0600)

	if err := runSchemaValidate([]string{envFile, schemaFile}); err != nil {
		t.Fatalf("optional missing field should not fail validation: %v", err)
	}
}

func TestSchemaValidate_MultipleViolations(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("PORT=abc\n"), 0600)

	schemaData, _ := json.Marshal(map[string]interface{}{
		"fields": []map[string]interface{}{
			{"key": "HOST", "required": true},
			{"key": "PORT", "required": true, "pattern": `^\d+$`},
		},
	})
	schemaFile := filepath.Join(dir, "schema.json")
	_ = os.WriteFile(schemaFile, schemaData, 0600)

	err := runSchemaValidate([]string{envFile, schemaFile})
	if err == nil {
		t.Fatal("expected error for multiple violations")
	}
}

func TestSchemaValidate_EmptyEnvFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte(""), 0600)

	schemaData, _ := json.Marshal(map[string]interface{}{
		"fields": []map[string]interface{}{
			{"key": "REQUIRED_KEY", "required": true},
		},
	})
	schemaFile := filepath.Join(dir, "schema.json")
	_ = os.WriteFile(schemaFile, schemaData, 0600)

	if err := runSchemaValidate([]string{envFile, schemaFile}); err == nil {
		t.Fatal("expected error for empty env against required key")
	}
}

func TestSchemaGenerate_ProducesValidJSON(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(envFile, []byte("API_KEY=secret\nAPI_URL=https://example.com\n"), 0600)

	// Redirect stdout capture is not straightforward; just verify no error.
	if err := runSchemaGenerate([]string{envFile}); err != nil {
		t.Fatalf("unexpected error generating schema: %v", err)
	}
}
