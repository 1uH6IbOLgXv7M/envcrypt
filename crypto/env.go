package crypto

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvFile represents a parsed .env file with its key-value pairs.
type EnvFile struct {
	Entries []EnvEntry
}

// EnvEntry holds a single line from a .env file, preserving comments and blanks.
type EnvEntry struct {
	Key     string
	Value   string
	Raw     string
	Comment bool
	Blank   bool
}

// ParseEnvFile reads and parses a .env file from the given path.
func ParseEnvFile(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	var ef EnvFile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		entry := parseLine(line)
		ef.Entries = append(ef.Entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}
	return &ef, nil
}

// Serialize converts the EnvFile back to its string representation.
func (ef *EnvFile) Serialize() string {
	var sb strings.Builder
	for _, e := range ef.Entries {
		sb.WriteString(e.Raw)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// ToMap returns a map of key-value pairs, skipping comments and blank lines.
func (ef *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		if !e.Comment && !e.Blank {
			m[e.Key] = e.Value
		}
	}
	return m
}

func parseLine(line string) EnvEntry {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return EnvEntry{Raw: line, Blank: true}
	}
	if strings.HasPrefix(trimmed, "#") {
		return EnvEntry{Raw: line, Comment: true}
	}
	parts := strings.SplitN(trimmed, "=", 2)
	if len(parts) != 2 {
		return EnvEntry{Raw: line, Blank: true}
	}
	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
	return EnvEntry{Key: key, Value: value, Raw: line}
}
