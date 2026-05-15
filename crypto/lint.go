package crypto

import "fmt"

// LintResult holds a single lint finding for an env file.
type LintResult struct {
	Line    int
	Key     string
	Message string
	Severity string // "warn" or "error"
}

// LintEnvFile parses an env file and returns a list of lint findings.
func LintEnvFile(path string) ([]LintResult, error) {
	entries, err := ParseEnvFile(path)
	if err != nil {
		return nil, fmt.Errorf("lint: %w", err)
	}

	var results []LintResult
	seen := make(map[string]int)

	for i, e := range entries {
		lineNum := i + 1

		// Duplicate key check
		if prev, ok := seen[e.Key]; ok {
			results = append(results, LintResult{
				Line:     lineNum,
				Key:      e.Key,
				Message:  fmt.Sprintf("duplicate key (first seen on line %d)", prev),
				Severity: "error",
			})
		} else {
			seen[e.Key] = lineNum
		}

		// Empty value warning
		if e.Value == "" {
			results = append(results, LintResult{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "empty value",
				Severity: "warn",
			})
		}

		// Key naming convention: should be UPPER_SNAKE_CASE
		if !isUpperSnakeCase(e.Key) {
			results = append(results, LintResult{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "key should be UPPER_SNAKE_CASE",
				Severity: "warn",
			})
		}
	}

	return results, nil
}

// isUpperSnakeCase returns true if s contains only uppercase letters, digits, and underscores.
func isUpperSnakeCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}
