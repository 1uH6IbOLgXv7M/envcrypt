package crypto

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField defines the expected shape of a single env variable.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Schema is a collection of field definitions.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation describes a single validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

// ValidateEnvMap checks an env map against the schema and returns violations.
func ValidateEnvMap(env map[string]string, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	for _, field := range schema.Fields {
		val, ok := env[field.Key]
		if !ok || strings.TrimSpace(val) == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}
		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				})
			}
		}
	}
	return violations
}

// FormatViolations returns a human-readable summary of violations.
func FormatViolations(violations []SchemaViolation) string {
	if len(violations) == 0 {
		return "schema validation passed: no violations found"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d schema violation(s):\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", v.Key, v.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}
