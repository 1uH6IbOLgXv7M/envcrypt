package crypto

import (
	"fmt"
)

// MultiRecipientEncrypt encrypts plaintext for multiple recipients.
// Each recipient gets an independently encrypted copy of the data,
// stored as a map of alias -> ciphertext.
func MultiRecipientEncrypt(plaintext []byte, recipients map[string]string) (map[string][]byte, error) {
	if len(recipients) == 0 {
		return nil, fmt.Errorf("no recipients provided")
	}

	result := make(map[string][]byte, len(recipients))
	for alias, pubKey := range recipients {
		ciphertext, err := Encrypt(plaintext, pubKey)
		if err != nil {
			return nil, fmt.Errorf("encrypt for recipient %q: %w", alias, err)
		}
		result[alias] = ciphertext
	}
	return result, nil
}

// MultiRecipientDecrypt attempts to decrypt a multi-recipient ciphertext map
// using the provided private key. It tries each entry until one succeeds.
func MultiRecipientDecrypt(ciphertexts map[string][]byte, privKey string) ([]byte, string, error) {
	if len(ciphertexts) == 0 {
		return nil, "", fmt.Errorf("no ciphertexts provided")
	}

	for alias, ct := range ciphertexts {
		plaintext, err := Decrypt(ct, privKey)
		if err == nil {
			return plaintext, alias, nil
		}
	}
	return nil, "", fmt.Errorf("could not decrypt with provided private key: no matching recipient found")
}
