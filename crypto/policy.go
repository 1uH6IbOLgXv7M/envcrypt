package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Policy defines rules for vault access and encryption requirements.
type Policy struct {
	MinRecipients   int      `json:"min_recipients"`
	RequireSign     bool     `json:"require_sign"`
	AllowedKeys     []string `json:"allowed_keys,omitempty"`
	MaxVersions     int      `json:"max_versions"`
	RequireAuditLog bool     `json:"require_audit_log"`
}

// PolicyPath returns the path to the policy file for a given directory.
func PolicyPath(dir string) string {
	return filepath.Join(dir, ".envcrypt", "policy.json")
}

// LoadPolicy loads the policy from the given directory.
// If no policy file exists, a default permissive policy is returned.
func LoadPolicy(dir string) (*Policy, error) {
	path := PolicyPath(dir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultPolicy(), nil
		}
		return nil, fmt.Errorf("load policy: %w", err)
	}
	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse policy: %w", err)
	}
	return &p, nil
}

// SavePolicy writes the policy to the given directory.
func SavePolicy(dir string, p *Policy) error {
	path := PolicyPath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("mkdir policy: %w", err)
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal policy: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// Validate checks whether the given vault and keyring satisfy the policy.
func (p *Policy) Validate(v *Vault, kr *KeyRing) error {
	if p.MinRecipients > 0 && len(kr.Keys) < p.MinRecipients {
		return fmt.Errorf("policy: need at least %d recipient(s), have %d", p.MinRecipients, len(kr.Keys))
	}
	if p.MaxVersions > 0 && len(v.Entries) > p.MaxVersions {
		return fmt.Errorf("policy: vault exceeds max versions %d (have %d)", p.MaxVersions, len(v.Entries))
	}
	if len(p.AllowedKeys) > 0 {
		allowed := make(map[string]bool, len(p.AllowedKeys))
		for _, k := range p.AllowedKeys {
			allowed[k] = true
		}
		for name := range kr.Keys {
			if !allowed[name] {
				return fmt.Errorf("policy: key %q is not in allowed_keys", name)
			}
		}
	}
	return nil
}

func defaultPolicy() *Policy {
	return &Policy{
		MinRecipients:   1,
		RequireSign:     false,
		MaxVersions:     100,
		RequireAuditLog: false,
	}
}
