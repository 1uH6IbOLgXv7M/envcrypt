package cmd

import (
	"fmt"
	"os"

	"envcrypt/crypto"
)

// runCompressInfo prints compression statistics for a given .env file.
// Usage: envcrypt compress-info <file>
func runCompressInfo(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt compress-info <env-file>")
	}
	filePath := args[0]

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("compress-info: cannot read file: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("File is empty.")
		return nil
	}

	compressed, err := crypto.CompressEnvMap(data)
	if err != nil {
		return fmt.Errorf("compress-info: compression failed: %w", err)
	}

	ratio, err := crypto.CompressionRatio(data)
	if err != nil {
		return fmt.Errorf("compress-info: ratio calculation failed: %w", err)
	}

	fmt.Printf("File:              %s\n", filePath)
	fmt.Printf("Original size:     %d bytes\n", len(data))
	fmt.Printf("Compressed size:   %d bytes\n", len(compressed))
	fmt.Printf("Compression ratio: %.2f\n", ratio)

	saved := len(data) - len(compressed)
	if saved > 0 {
		fmt.Printf("Space saved:       %d bytes (%.1f%%)\n", saved, (1.0-ratio)*100)
	} else {
		fmt.Println("Note: file is too small to benefit from compression.")
	}
	return nil
}
