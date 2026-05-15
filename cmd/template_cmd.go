package cmd

import (
	"fmt"
	"os"
	"strings"

	"envcrypt/crypto"
)

// runTemplateRender renders a template file against a decrypted vault entry.
// Usage: envcrypt template render <template-file> [--version N] [--output file]
func runTemplateRender(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: template render <template-file> [--version N] [--output <file>]")
	}

	tmplPath := args[0]
	version := -1
	outputPath := ""

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--version":
			if i+1 >= len(args) {
				return fmt.Errorf("--version requires a value")
			}
			i++
			fmt.Sscanf(args[i], "%d", &version)
		case "--output":
			if i+1 >= len(args) {
				return fmt.Errorf("--output requires a value")
			}
			i++
			outputPath = args[i]
		}
	}

	privKey, err := crypto.LoadPrivateKey(".envcrypt/private.key")
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	vault, err := crypto.LoadVault(".envcrypt/vault.json")
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	var entry *crypto.VaultEntry
	if version >= 0 {
		e, ok := crypto.EntryByVersion(vault, version)
		if !ok {
			return fmt.Errorf("version %d not found in vault", version)
		}
		entry = &e
	} else {
		e, ok := crypto.LatestEntry(vault)
		if !ok {
			return fmt.Errorf("vault is empty")
		}
		entry = &e
	}

	plaintext, err := crypto.Decrypt(privKey, entry.Ciphertext)
	if err != nil {
		return fmt.Errorf("decrypt vault entry: %w", err)
	}

	env, err := crypto.ParseEnvBytes(plaintext)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	result, err := crypto.RenderTemplateFile(tmplPath, env)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if len(result.Missing) > 0 {
		fmt.Fprintf(os.Stderr, "warning: %s\n", crypto.FormatMissing(result.Missing))
	}

	output := strings.TrimRight(result.Rendered, "\n") + "\n"
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(output), 0o600); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Printf("rendered template written to %s\n", outputPath)
		return nil
	}

	fmt.Print(output)
	return nil
}
