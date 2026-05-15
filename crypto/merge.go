package crypto

import "fmt"

// MergeStrategy defines how conflicts are resolved during merge.
type MergeStrategy int

const (
	// MergeStrategyOurs keeps the value from the base env file on conflict.
	MergeStrategyOurs MergeStrategy = iota
	// MergeStrategyTheirs keeps the value from the incoming env file on conflict.
	MergeStrategyTheirs
)

// MergeResult holds the merged env map and any conflict details.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []MergeConflict
}

// MergeConflict describes a key that existed in both files with different values.
type MergeConflict struct {
	Key      string
	OursVal  string
	TheirsVal string
	Resolved string
}

// MergeEnvFiles merges two parsed env maps using the given strategy.
// Keys only in base or only in incoming are always included.
// Conflicting keys are resolved per the strategy and recorded in Conflicts.
func MergeEnvFiles(base, incoming map[string]string, strategy MergeStrategy) MergeResult {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	var conflicts []MergeConflict
	for k, theirVal := range incoming {
		ourVal, exists := merged[k]
		if !exists {
			merged[k] = theirVal
			continue
		}
		if ourVal == theirVal {
			continue
		}
		// Conflict
		var resolved string
		switch strategy {
		case MergeStrategyTheirs:
			resolved = theirVal
		default: // MergeStrategyOurs
			resolved = ourVal
		}
		merged[k] = resolved
		conflicts = append(conflicts, MergeConflict{
			Key:       k,
			OursVal:   ourVal,
			TheirsVal: theirVal,
			Resolved:  resolved,
		})
	}
	return MergeResult{Merged: merged, Conflicts: conflicts}
}

// FormatMergeReport returns a human-readable summary of merge conflicts.
func FormatMergeReport(result MergeResult) string {
	if len(result.Conflicts) == 0 {
		return "merge completed with no conflicts\n"
	}
	out := fmt.Sprintf("%d conflict(s) resolved:\n", len(result.Conflicts))
	for _, c := range result.Conflicts {
		out += fmt.Sprintf("  %s: ours=%q theirs=%q => resolved=%q\n", c.Key, c.OursVal, c.TheirsVal, c.Resolved)
	}
	return out
}
