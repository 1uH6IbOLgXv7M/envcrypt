package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envcrypt/crypto"
)

func setupMergeTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "merge-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeMergeEnv(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
}

func TestRunMerge_NoArgs(t *testing.T) {
	// Should not panic; exits with error. We test the logic layer instead.
	base := map[string]string{"A": "1"}
	incoming := map[string]string{"B": "2"}
	result := crypto.MergeEnvFiles(base, incoming, crypto.MergeStrategyOurs)
	if len(result.Merged) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result.Merged))
	}
}

func TestRunMerge_RoundTrip_NoConflict(t *testing.T) {
	dir := setupMergeTestDir(t)
	basePath := filepath.Join(dir, "base.env")
	incomingPath := filepath.Join(dir, "incoming.env")
	outPath := filepath.Join(dir, "merged.env")

	writeMergeEnv(t, basePath, "A=1\nB=2\n")
	writeMergeEnv(t, incomingPath, "C=3\n")

	baseMap, _ := crypto.ParseEnvFile(basePath)
	incomingMap, _ := crypto.ParseEnvFile(incomingPath)
	result := crypto.MergeEnvFiles(baseMap, incomingMap, crypto.MergeStrategyOurs)

	env := crypto.EnvFile{}
	for k, v := range result.Merged {
		env.Entries = append(env.Entries, crypto.EnvEntry{Key: k, Value: v})
	}
	if err := os.WriteFile(outPath, []byte(env.Serialize()), 0600); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	mergedMap, err := crypto.ParseEnvFile(outPath)
	if err != nil {
		t.Fatalf("parse merged file failed: %v", err)
	}
	if mergedMap["A"] != "1" || mergedMap["B"] != "2" || mergedMap["C"] != "3" {
		t.Errorf("merged map incorrect: %v", mergedMap)
	}
}

func TestRunMerge_ConflictReport(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}
	result := crypto.MergeEnvFiles(base, incoming, crypto.MergeStrategyTheirs)
	report := crypto.FormatMergeReport(result)
	if !strings.Contains(report, "KEY") {
		t.Errorf("report should mention conflicted key, got: %q", report)
	}
	if result.Merged["KEY"] != "new" {
		t.Errorf("expected theirs strategy to resolve to 'new'")
	}
}
