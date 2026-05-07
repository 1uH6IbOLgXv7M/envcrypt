package crypto

import (
	"fmt"
	"os"
	"strings"
)

// Recipient represents a named public key that can receive encrypted secrets.
type Recipient struct {
	Alias     string `json:"alias"`
	PublicKey string `json:"public_key"`
}

// RecipientsFile holds a list of recipients for multi-recipient encryption.
type RecipientsFile struct {
	Recipients []Recipient `json:"recipients"`
}

// LoadRecipientsFile reads a .env-recipients file from disk.
func LoadRecipientsFile(path string) (*RecipientsFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &RecipientsFile{}, nil
		}
		return nil, fmt.Errorf("read recipients file: %w", err)
	}
	rf := &RecipientsFile{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed recipients line: %q", line)
		}
		rf.Recipients = append(rf.Recipients, Recipient{
			Alias:     strings.TrimSpace(parts[0]),
			PublicKey: strings.TrimSpace(parts[1]),
		})
	}
	return rf, nil
}

// SaveRecipientsFile writes the recipients list to disk.
func SaveRecipientsFile(path string, rf *RecipientsFile) error {
	var sb strings.Builder
	sb.WriteString("# envcrypt recipients file — add public keys for multi-recipient encryption\n")
	for _, r := range rf.Recipients {
		sb.WriteString(fmt.Sprintf("%s=%s\n", r.Alias, r.PublicKey))
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// AddRecipient adds or overwrites a recipient by alias.
func (rf *RecipientsFile) AddRecipient(alias, pubKey string) {
	for i, r := range rf.Recipients {
		if r.Alias == alias {
			rf.Recipients[i].PublicKey = pubKey
			return
		}
	}
	rf.Recipients = append(rf.Recipients, Recipient{Alias: alias, PublicKey: pubKey})
}

// GetRecipient returns the public key for the given alias, or an error if not found.
func (rf *RecipientsFile) GetRecipient(alias string) (string, error) {
	for _, r := range rf.Recipients {
		if r.Alias == alias {
			return r.PublicKey, nil
		}
	}
	return "", fmt.Errorf("recipient %q not found", alias)
}

// RemoveRecipient removes a recipient by alias. Returns false if not found.
func (rf *RecipientsFile) RemoveRecipient(alias string) bool {
	for i, r := range rf.Recipients {
		if r.Alias == alias {
			rf.Recipients = append(rf.Recipients[:i], rf.Recipients[i+1:]...)
			return true
		}
	}
	return false
}
