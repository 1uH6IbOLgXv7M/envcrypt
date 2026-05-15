package cmd

import (
	"fmt"
	"os"

	"envcrypt/crypto"
)

// runLint lints an env file for common issues.
// Usage: envcrypt lint <file>
func runLint(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envcrypt lint <file>")
		os.Exit(1)
	}

	path := args[0]
	results, err := crypto.LintEnvFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lint error: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Printf("✔  %s — no issues found\n", path)
		return
	}

	errorCount := 0
	warnCount := 0
	for _, r := range results {
		switch r.Severity {
		case "error":
			errorCount++
			fmt.Printf("[ERROR] line %d  %-24s %s\n", r.Line, r.Key, r.Message)
		case "warn":
			warnCount++
			fmt.Printf("[WARN]  line %d  %-24s %s\n", r.Line, r.Key, r.Message)
		}
	}

	fmt.Printf("\n%d error(s), %d warning(s) in %s\n", errorCount, warnCount, path)

	if errorCount > 0 {
		os.Exit(2)
	}
}
