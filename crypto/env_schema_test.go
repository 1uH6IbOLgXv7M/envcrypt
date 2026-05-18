package crypto

import (
	"strings"
	"testing"
)

func TestValidateEnvMap_AllPresent(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	schema := Schema{Fields: []SchemaField{
		{Key: "HOST", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	}}
	violations := ValidateEnvMap(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestValidateEnvMap_MissingRequired(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	schema := Schema{Fields: []SchemaField{
		{Key: "HOST", Required: true},
		{Key: "PORT", Required: true},
	}}
	violations := ValidateEnvMap(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %s", violations[0].Key)
	}
}

func TestValidateEnvMap_PatternMismatch(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	schema := Schema{Fields: []SchemaField{
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	}}
	violations := ValidateEnvMap(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "does not match pattern") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestValidateEnvMap_OptionalMissing(t *testing.T) {
	env := map[string]string{}
	schema := Schema{Fields: []SchemaField{
		{Key: "DEBUG", Required: false},
	}}
	violations := ValidateEnvMap(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for optional missing key, got %d", len(violations))
	}
}

func TestValidateEnvMap_InvalidPattern(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	schema := Schema{Fields: []SchemaField{
		{Key: "KEY", Required: true, Pattern: `[invalid`},
	}}
	violations := ValidateEnvMap(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for invalid pattern, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "invalid pattern") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestFormatViolations_NoViolations(t *testing.T) {
	out := FormatViolations(nil)
	if !strings.Contains(out, "passed") {
		t.Errorf("expected pass message, got: %s", out)
	}
}

func TestFormatViolations_WithViolations(t *testing.T) {
	v := []SchemaViolation{{Key: "HOST", Message: "required key is missing or empty"}}
	out := FormatViolations(v)
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "1 schema violation") {
		t.Errorf("expected violation count in output, got: %s", out)
	}
}
