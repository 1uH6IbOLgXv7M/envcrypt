package crypto

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderTemplate_Basic(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	res := RenderTemplate("postgres://${HOST}:${PORT}/db", env)
	want := "postgres://localhost:5432/db"
	if res.Rendered != want {
		t.Errorf("got %q, want %q", res.Rendered, want)
	}
	if len(res.Missing) != 0 {
		t.Errorf("unexpected missing keys: %v", res.Missing)
	}
}

func TestRenderTemplate_MissingKeys(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	res := RenderTemplate("${HOST}:${PORT}", env)
	if !strings.Contains(res.Rendered, "${PORT}") {
		t.Errorf("expected placeholder to remain, got %q", res.Rendered)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "PORT" {
		t.Errorf("expected missing PORT, got %v", res.Missing)
	}
}

func TestRenderTemplate_DuplicateMissingKey(t *testing.T) {
	res := RenderTemplate("${X} and ${X} again", map[string]string{})
	if len(res.Missing) != 1 {
		t.Errorf("expected deduplicated missing key, got %v", res.Missing)
	}
}

func TestRenderTemplate_NoPlaceholders(t *testing.T) {
	res := RenderTemplate("no placeholders here", map[string]string{})
	if res.Rendered != "no placeholders here" {
		t.Errorf("unexpected change: %q", res.Rendered)
	}
	if len(res.Missing) != 0 {
		t.Errorf("unexpected missing: %v", res.Missing)
	}
}

func TestRenderTemplateFile_Basic(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.txt")
	_ = os.WriteFile(p, []byte("hello ${NAME}!"), 0o644)
	res, err := RenderTemplateFile(p, map[string]string{"NAME": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "hello world!" {
		t.Errorf("got %q", res.Rendered)
	}
}

func TestRenderTemplateFile_Missing(t *testing.T) {
	_, err := RenderTemplateFile("/nonexistent/path/tmpl.txt", nil)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFormatMissing(t *testing.T) {
	if s := FormatMissing(nil); s != "" {
		t.Errorf("expected empty, got %q", s)
	}
	s := FormatMissing([]string{"A", "B"})
	if !strings.Contains(s, "A") || !strings.Contains(s, "B") {
		t.Errorf("unexpected format: %q", s)
	}
}
