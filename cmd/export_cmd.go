package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"envcrypt/crypto"
)

// runExport handles the `envcrypt export` subcommand.
// Usage: envcrypt export [--version N] [--format json|dotenv|shell] [--out FILE] <vault>
func runExport(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt export [--version N] [--format json|dotenv|shell] [--out FILE] <vault>")
	}

	vaultPath := args[0]
	format := crypto.FormatDotEnv
	version := 0
	outPath := "-"

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--format":
			if i+1 >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			i++
			format = crypto.ExportFormat(args[i])
		case "--version":
			if i+1 >= len(args) {
				return fmt.Errorf("--version requires a value")
			}
			i++
			v, err := strconv.Atoi(args[i])
			if err != nil {
				return fmt.Errorf("invalid version: %s", args[i])
			}
			version = v
		case "--out":
			if i+1 >= len(args) {
				return fmt.Errorf("--out requires a value")
			}
			i++
			outPath = args[i]
		}
	}

	privKeyPath := filepath.Join(filepath.Dir(vaultPath), "private.age")
	if envKey := os.Getenv("ENVCRYPT_PRIVATE_KEY"); envKey != "" {
		privKeyPath = envKey
	}

	entry, err := crypto.ExportVault(vaultPath, privKeyPath, version)
	if err != nil {
		return fmt.Errorf("export vault: %w", err)
	}

	output, err := crypto.RenderExport(entry, format)
	if err != nil {
		return fmt.Errorf("render export: %w", err)
	}

	if err := crypto.WriteExport(output, outPath); err != nil {
		return fmt.Errorf("write export: %w", err)
	}

	if outPath != "-" {
		fmt.Fprintf(os.Stderr, "exported version %d to %s\n", entry.Version, outPath)
	}
	return nil
}
