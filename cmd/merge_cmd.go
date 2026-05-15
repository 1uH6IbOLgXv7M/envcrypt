package cmd

import (
	"fmt"
	"os"

	"envcrypt/crypto"
)

// runMerge merges two .env files and writes the result.
// Usage: envcrypt merge <base.env> <incoming.env> [--strategy=ours|theirs] [--output=out.env]
func runMerge(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envcrypt merge <base.env> <incoming.env> [--strategy=ours|theirs] [--output=<file>]")
		os.Exit(1)
	}

	basePath := args[0]
	incomingPath := args[1]

	strategy := crypto.MergeStrategyOurs
	outputPath := ""

	for _, arg := range args[2:] {
		switch {
		case arg == "--strategy=theirs":
			strategy = crypto.MergeStrategyTheirs
		case arg == "--strategy=ours":
			strategy = crypto.MergeStrategyOurs
		case len(arg) > 9 && arg[:9] == "--output=":
			outputPath = arg[9:]
		}
	}

	baseMap, err := crypto.ParseEnvFile(basePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file: %v\n", err)
		os.Exit(1)
	}

	incomingMap, err := crypto.ParseEnvFile(incomingPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading incoming file: %v\n", err)
		os.Exit(1)
	}

	result := crypto.MergeEnvFiles(baseMap, incomingMap, strategy)
	report := crypto.FormatMergeReport(result)
	fmt.Print(report)

	env := crypto.EnvFile{Entries: make([]crypto.EnvEntry, 0, len(result.Merged))}
	for k, v := range result.Merged {
		env.Entries = append(env.Entries, crypto.EnvEntry{Key: k, Value: v})
	}

	if outputPath == "" {
		fmt.Print(env.Serialize())
		return
	}

	if err := os.WriteFile(outputPath, []byte(env.Serialize()), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("merged output written to %s\n", outputPath)
}
