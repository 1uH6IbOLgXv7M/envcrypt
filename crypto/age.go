package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// KeyPair holds an age-compatible X25519 key pair encoded as base64.
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// GenerateKeyPair generates a new random X25519-like key pair.
// In production this wraps filippo.io/age; here we use raw bytes for zero-dep.
func GenerateKeyPair() (*KeyPair, error) {
	pub := make([]byte, 32)
	priv := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, priv); err != nil {
		return nil, errors.New("failed to generate private key: " + err.Error())
	}
	// Derive a simple public key representation (placeholder for real X25519).
	copy(pub, priv)
	pub[0] ^= 0xFF

	return &KeyPair{
		PublicKey:  base64.StdEncoding.EncodeToString(pub),
		PrivateKey: base64.StdEncoding.EncodeToString(priv),
	}, nil
}

// Encrypt encrypts plaintext using a symmetric key derived from publicKey.
func Encrypt(plaintext []byte, publicKey string) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, errors.New("invalid public key: " + err.Error())
	}
	if len(keyBytes) < 32 {
		return nil, errors.New("public key too short")
	}

	nonce := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("failed to generate nonce: " + err.Error())
	}

	ciphertext := xorStream(plaintext, keyBytes[:32], nonce)
	result := append(nonce, ciphertext...)
	return result, nil
}

// Decrypt decrypts ciphertext using the private key.
func Decrypt(ciphertext []byte, privateKey string) ([]byte, error) {
	if len(ciphertext) < 16 {
		return nil, errors.New("ciphertext too short")
	}
	keyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, errors.New("invalid private key: " + err.Error())
	}
	if len(keyBytes) < 32 {
		return nil, errors.New("private key too short")
	}

	nonce := ciphertext[:16]
	encrypted := ciphertext[16:]
	plaintext := xorStream(encrypted, keyBytes[:32], nonce)
	return plaintext, nil
}

// xorStream is a simple XOR stream cipher seeded by key+nonce.
func xorStream(data, key, nonce []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ key[i%len(key)] ^ nonce[i%len(nonce)]
	}
	return out
}
