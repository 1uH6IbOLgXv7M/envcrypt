package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExportFormat defines the output format for exported env entries.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatDotEnv ExportFormat = "dotenv"
	FormatShell ExportFormat = "shell"
)

// ExportEntry holds a single decrypted env entry for export.
type ExportEntry struct {
	Version   int               `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Vars      map[string]string `json:"vars"`
}

// ExportVault decrypts the latest (or specified) vault version and returns an ExportEntry.
func ExportVault(vaultPath, privateKeyPath string, version int) (*ExportEntry, error) {
	vault, err := LoadVault(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	var entry *VaultEntry
	if version <= 0 {
		entry = vault.Latest()
	} else {
		entry = vault.ByVersion(version)
	}
	if entry == nil {
		return nil, fmt.Errorf("vault entry not found")
	}

	privKey, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load private key: %w", err)
	}

	plaintext, err := Decrypt(privKey, entry.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	envFile, err := ParseEnvBytes(plaintext)
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	vars := make(map[string]string, len(envFile.Entries))
	for _, e := range envFile.Entries {
		vars[e.Key] = e.Value
	}

	return &ExportEntry{
		Version:   entry.Version,
		Timestamp: entry.CreatedAt,
		Vars:      vars,
	}, nil
}

// RenderExport formats an ExportEntry into the requested format.
func RenderExport(entry *ExportEntry, format ExportFormat) (string, error) {
	switch format {
	case FormatJSON:
		b, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	case FormatDotEnv:
		out := ""
		for k, v := range entry.Vars {
			out += fmt.Sprintf("%s=%s\n", k, v)
		}
		return out, nil
	case FormatShell:
		out := ""
		for k, v := range entry.Vars {
			out += fmt.Sprintf("export %s=%q\n", k, v)
		}
		return out, nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

// WriteExport writes rendered export output to a file (or stdout if path is "-").
func WriteExport(output, path string) error {
	if path == "-" {
		_, err := fmt.Print(output)
		return err
	}
	return os.WriteFile(path, []byte(output), 0600)
}
