package crypto

import (
	"bufio"
	"bytes"
	"fmt"
)

// ParseEnvBytes parses a .env file from raw bytes instead of a file path.
// It reuses the same line-parsing logic as ParseEnvFile.
func ParseEnvBytes(data []byte) (*EnvFile, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty env data")
	}

	ef := &EnvFile{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		entry, ok, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
		if ok {
			ef.Entries = append(ef.Entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return ef, nil
}
