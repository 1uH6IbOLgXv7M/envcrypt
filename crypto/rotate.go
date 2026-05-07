package crypto

import "fmt"

// RotateResult holds the outcome of a key rotation operation.
type RotateResult struct {
	EntriesRotated int
	NewPublicKey   string
}

// RotateVaultKeys re-encrypts all vault entries with a new key pair.
// It loads the vault at vaultPath, decrypts each entry with oldPrivKey,
// re-encrypts with newPubKey, and saves the updated vault.
func RotateVaultKeys(vaultPath, oldPrivKey, newPubKey string) (*RotateResult, error) {
	vault, err := LoadVault(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("rotate: load vault: %w", err)
	}

	if len(vault.Entries) == 0 {
		return &RotateResult{EntriesRotated: 0, NewPublicKey: newPubKey}, nil
	}

	for i, entry := range vault.Entries {
		plaintext, err := Decrypt(oldPrivKey, entry.Ciphertext)
		if err != nil {
			return nil, fmt.Errorf("rotate: decrypt entry version %d: %w", entry.Version, err)
		}

		newCiphertext, err := Encrypt(newPubKey, plaintext)
		if err != nil {
			return nil, fmt.Errorf("rotate: encrypt entry version %d: %w", entry.Version, err)
		}

		vault.Entries[i].Ciphertext = newCiphertext
	}

	if err := SaveVault(vaultPath, vault); err != nil {
		return nil, fmt.Errorf("rotate: save vault: %w", err)
	}

	return &RotateResult{
		EntriesRotated: len(vault.Entries),
		NewPublicKey:   newPubKey,
	}, nil
}
