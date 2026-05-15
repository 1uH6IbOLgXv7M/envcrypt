package crypto

import (
	"strings"
	"testing"
)

func TestMergeEnvFiles_NoConflict(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"C": "3"}
	result := MergeEnvFiles(base, incoming, MergeStrategyOurs)
	if len(result.Conflicts) != 0 {
		t.Fatalf("expected 0 conflicts, got %d", len(result.Conflicts))
	}
	if result.Merged["A"] != "1" || result.Merged["B"] != "2" || result.Merged["C"] != "3" {
		t.Errorf("unexpected merged map: %v", result.Merged)
	}
}

func TestMergeEnvFiles_ConflictOurs(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	incoming := map[string]string{"KEY": "their_val"}
	result := MergeEnvFiles(base, incoming, MergeStrategyOurs)
	if len(result.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(result.Conflicts))
	}
	if result.Merged["KEY"] != "base_val" {
		t.Errorf("expected ours strategy to keep base_val, got %q", result.Merged["KEY"])
	}
	if result.Conflicts[0].Resolved != "base_val" {
		t.Errorf("conflict resolved value mismatch")
	}
}

func TestMergeEnvFiles_ConflictTheirs(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	incoming := map[string]string{"KEY": "their_val"}
	result := MergeEnvFiles(base, incoming, MergeStrategyTheirs)
	if result.Merged["KEY"] != "their_val" {
		t.Errorf("expected theirs strategy to keep their_val, got %q", result.Merged["KEY"])
	}
}

func TestMergeEnvFiles_SameValue_NoConflict(t *testing.T) {
	base := map[string]string{"KEY": "same"}
	incoming := map[string]string{"KEY": "same"}
	result := MergeEnvFiles(base, incoming, MergeStrategyOurs)
	if len(result.Conflicts) != 0 {
		t.Errorf("identical values should not be a conflict")
	}
}

func TestFormatMergeReport_NoConflicts(t *testing.T) {
	result := MergeResult{Merged: map[string]string{"A": "1"}}
	report := FormatMergeReport(result)
	if !strings.Contains(report, "no conflicts") {
		t.Errorf("expected no-conflict message, got: %q", report)
	}
}

func TestFormatMergeReport_WithConflicts(t *testing.T) {
	result := MergeResult{
		Merged: map[string]string{"X": "a"},
		Conflicts: []MergeConflict{
			{Key: "X", OursVal: "a", TheirsVal: "b", Resolved: "a"},
		},
	}
	report := FormatMergeReport(result)
	if !strings.Contains(report, "1 conflict") {
		t.Errorf("expected conflict count in report, got: %q", report)
	}
	if !strings.Contains(report, "X") {
		t.Errorf("expected key name in report")
	}
}
