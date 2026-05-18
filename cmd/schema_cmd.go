package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"envcrypt/crypto"
)

// runSchemaValidate validates a .env file against a JSON schema file.
// Usage: envcrypt schema validate <env-file> <schema-file>
func runSchemaValidate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: schema validate <env-file> <schema-file>")
	}
	envFile, schemaFile := args[0], args[1]

	env, err := crypto.ParseEnvFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to parse env file: %w", err)
	}

	schema, err := loadSchemaFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to load schema file: %w", err)
	}

	violations := crypto.ValidateEnvMap(env, schema)
	fmt.Println(crypto.FormatViolations(violations))
	if len(violations) > 0 {
		return fmt.Errorf("schema validation failed with %d violation(s)", len(violations))
	}
	return nil
}

// runSchemaGenerate generates a starter schema JSON from an existing .env file.
// Usage: envcrypt schema generate <env-file>
func runSchemaGenerate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: schema generate <env-file>")
	}
	env, err := crypto.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse env file: %w", err)
	}

	type jsonField struct {
		Key      string `json:"key"`
		Required bool   `json:"required"`
		Pattern  string `json:"pattern,omitempty"`
	}
	var fields []jsonField
	for k := range env {
		fields = append(fields, jsonField{Key: k, Required: true})
	}

	out, err := json.MarshalIndent(map[string]interface{}{"fields": fields}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}
	fmt.Println(string(out))
	return nil
}

// loadSchemaFile reads a JSON schema file and returns a Schema.
func loadSchemaFile(path string) (crypto.Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return crypto.Schema{}, err
	}
	type jsonField struct {
		Key      string `json:"key"`
		Required string `json:"required"` // accept bool or string
		Pattern  string `json:"pattern"`
	}
	var raw struct {
		Fields []json.RawMessage `json:"fields"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return crypto.Schema{}, fmt.Errorf("invalid schema JSON: %w", err)
	}
	var schema crypto.Schema
	for _, rf := range raw.Fields {
		var m map[string]interface{}
		if err := json.Unmarshal(rf, &m); err != nil {
			return crypto.Schema{}, err
		}
		field := crypto.SchemaField{}
		if v, ok := m["key"].(string); ok {
			field.Key = v
		}
		if v, ok := m["pattern"].(string); ok {
			field.Pattern = v
		}
		switch r := m["required"].(type) {
		case bool:
			field.Required = r
		case string:
			field.Required, _ = strconv.ParseBool(strings.ToLower(r))
		}
		schema.Fields = append(schema.Fields, field)
	}
	return schema, nil
}
