package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// EnvChecksum holds the SHA-256 checksum of a set of env key=value pairs.
type EnvChecksum struct {
	Hash    string            `json:"hash"`
	Keys    []string          `json:"keys"`
	Entries map[string]string `json:"entries"`
}

// ChecksumEnvMap computes a deterministic SHA-256 checksum over the given
// key/value map. Keys are sorted before hashing to ensure stability.
func ChecksumEnvMap(env map[string]string) EnvChecksum {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return EnvChecksum{
		Hash:    hex.EncodeToString(sum[:]),
		Keys:    keys,
		Entries: env,
	}
}

// ChecksumEnvFile parses the given file and returns its checksum.
func ChecksumEnvFile(path string) (EnvChecksum, error) {
	env, err := ParseEnvFile(path)
	if err != nil {
		return EnvChecksum{}, fmt.Errorf("checksum: parse env file: %w", err)
	}
	return ChecksumEnvMap(env), nil
}

// MatchChecksum returns true when the provided hex hash matches the checksum
// of the given env map.
func MatchChecksum(env map[string]string, expected string) bool {
	c := ChecksumEnvMap(env)
	return c.Hash == expected
}
