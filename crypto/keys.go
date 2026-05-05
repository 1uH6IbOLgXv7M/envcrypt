package crypto

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultKeyDir      = ".envcrypt"
	PublicKeyFile      = "key.pub"
	PrivateKeyFile     = "key"
	PublicKeyPrefix    = "# envcrypt public key\n"
	PrivateKeyPrefix   = "# envcrypt private key — keep secret!\n"
)

// SaveKeyPair writes the key pair to dir (default: ~/.envcrypt/).
func SaveKeyPair(kp *KeyPair, dir string) error {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		dir = filepath.Join(home, DefaultKeyDir)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create key directory: %w", err)
	}

	pubPath := filepath.Join(dir, PublicKeyFile)
	privPath := filepath.Join(dir, PrivateKeyFile)

	if err := os.WriteFile(pubPath, []byte(PublicKeyPrefix+kp.PublicKey+"\n"), 0644); err != nil {
		return fmt.Errorf("cannot write public key: %w", err)
	}
	if err := os.WriteFile(privPath, []byte(PrivateKeyPrefix+kp.PrivateKey+"\n"), 0600); err != nil {
		return fmt.Errorf("cannot write private key: %w", err)
	}
	return nil
}

// LoadPublicKey reads the public key from a file, stripping comment lines.
func LoadPublicKey(path string) (string, error) {
	return loadKey(path)
}

// LoadPrivateKey reads the private key from a file, stripping comment lines.
func LoadPrivateKey(path string) (string, error) {
	return loadKey(path)
}

func loadKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read key file %q: %w", path, err)
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return "", errors.New("key file is empty or contains only comments")
	}
	return lines[0], nil
}
