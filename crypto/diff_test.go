package crypto

import (
	"strings"
	"testing"
)

func TestDiffEnvFiles_Added(t *testing.T) {
	old := EnvFile{"A": "1"}
	new := EnvFile{"A": "1", "B": "2"}
	d := DiffEnvFiles(old, new)
	if _, ok := d.Added["B"]; !ok {
		t.Error("expected B to be in Added")
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Error("unexpected diff entries")
	}
}

func TestDiffEnvFiles_Removed(t *testing.T) {
	old := EnvFile{"A": "1", "B": "2"}
	new := EnvFile{"A": "1"}
	d := DiffEnvFiles(old, new)
	if _, ok := d.Removed["B"]; !ok {
		t.Error("expected B to be in Removed")
	}
}

func TestDiffEnvFiles_Changed(t *testing.T) {
	old := EnvFile{"A": "1"}
	new := EnvFile{"A": "2"}
	d := DiffEnvFiles(old, new)
	v, ok := d.Changed["A"]
	if !ok {
		t.Fatal("expected A to be in Changed")
	}
	if v[0] != "1" || v[1] != "2" {
		t.Errorf("unexpected changed values: %v", v)
	}
}

func TestDiffEnvFiles_Unchanged(t *testing.T) {
	old := EnvFile{"A": "1"}
	new := EnvFile{"A": "1"}
	d := DiffEnvFiles(old, new)
	if _, ok := d.Unchanged["A"]; !ok {
		t.Error("expected A to be in Unchanged")
	}
	if len(d.Added)+len(d.Removed)+len(d.Changed) != 0 {
		t.Error("unexpected diff entries")
	}
}

func TestDiffEnvFiles_Empty(t *testing.T) {
	d := DiffEnvFiles(EnvFile{}, EnvFile{})
	if len(d.Added)+len(d.Removed)+len(d.Changed)+len(d.Unchanged) != 0 {
		t.Error("expected empty diff for two empty files")
	}
}

func TestFormatDiff_ContainsSymbols(t *testing.T) {
	d := DiffResult{
		Added:     map[string]string{"NEW": "val"},
		Removed:   map[string]string{"OLD": "val"},
		Changed:   map[string][2]string{"X": {"a", "b"}},
		Unchanged: map[string]string{},
	}
	out := FormatDiff(d)
	if !strings.Contains(out, "+ NEW") {
		t.Error("expected added line")
	}
	if !strings.Contains(out, "- OLD") {
		t.Error("expected removed line")
	}
	if !strings.Contains(out, "~ X") {
		t.Error("expected changed line")
	}
}

func TestFormatDiff_Empty(t *testing.T) {
	d := DiffResult{
		Added:     map[string]string{},
		Removed:   map[string]string{},
		Changed:   map[string][2]string{},
		Unchanged: map[string]string{},
	}
	out := FormatDiff(d)
	if strings.ContainsAny(out, "+-~") {
		t.Errorf("expected no diff symbols for empty diff, got: %q", out)
	}
}
