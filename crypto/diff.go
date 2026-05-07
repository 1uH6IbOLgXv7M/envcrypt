package crypto

import "fmt"

// DiffResult holds the comparison between two env file versions.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
	Unchanged map[string]string
}

// DiffEnvFiles compares two EnvFile maps and returns a DiffResult.
func DiffEnvFiles(oldEnv, newEnv EnvFile) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for key, newVal := range newEnv {
		if oldVal, exists := oldEnv[key]; exists {
			if oldVal != newVal {
				result.Changed[key] = [2]string{oldVal, newVal}
			} else {
				result.Unchanged[key] = newVal
			}
		} else {
			result.Added[key] = newVal
		}
	}

	for key, oldVal := range oldEnv {
		if _, exists := newEnv[key]; !exists {
			result.Removed[key] = oldVal
		}
	}

	return result
}

// FormatDiff returns a human-readable string representation of a DiffResult.
func FormatDiff(d DiffResult) string {
	out := ""
	for k, v := range d.Added {
		out += fmt.Sprintf("+ %s=%s\n", k, v)
	}
	for k, v := range d.Removed {
		out += fmt.Sprintf("- %s=%s\n", k, v)
	}
	for k, v := range d.Changed {
		out += fmt.Sprintf("~ %s: %s -> %s\n", k, v[0], v[1])
	}
	return out
}
