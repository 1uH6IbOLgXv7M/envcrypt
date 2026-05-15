package crypto

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateResult holds the rendered output and any missing keys.
type TemplateResult struct {
	Rendered string
	Missing  []string
}

var templateVarRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// RenderTemplate substitutes ${KEY} placeholders in the template string
// with values from the provided env map. Missing keys are collected rather
// than causing an error, allowing callers to decide how to handle them.
func RenderTemplate(tmpl string, env map[string]string) TemplateResult {
	seen := map[string]bool{}
	var missing []string

	rendered := templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := templateVarRe.FindStringSubmatch(match)[1]
		if val, ok := env[key]; ok {
			return val
		}
		if !seen[key] {
			seen[key] = true
			missing = append(missing, key)
		}
		return match // leave placeholder intact
	})

	return TemplateResult{Rendered: rendered, Missing: missing}
}

// RenderTemplateFile reads a template file from disk, substitutes placeholders
// using the provided env map, and returns the result.
func RenderTemplateFile(path string, env map[string]string) (TemplateResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TemplateResult{}, fmt.Errorf("read template file: %w", err)
	}
	return RenderTemplate(string(data), env), nil
}

// FormatMissing returns a human-readable summary of missing keys.
func FormatMissing(missing []string) string {
	if len(missing) == 0 {
		return ""
	}
	return fmt.Sprintf("missing keys: %s", strings.Join(missing, ", "))
}
