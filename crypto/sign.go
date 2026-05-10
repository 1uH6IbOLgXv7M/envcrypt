package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// SignatureKeyPair holds an Ed25519 signing key pair.
type SignatureKeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// GenerateSigningKeyPair generates a new Ed25519 key pair for signing.
func GenerateSigningKeyPair() (*SignatureKeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &SignatureKeyPair{PublicKey: pub, PrivateKey: priv}, nil
}

// SignVault produces a hex-encoded Ed25519 signature over the SHA-256 hash
// of the provided vault ciphertext blob.
func SignVault(priv ed25519.PrivateKey, ciphertext []byte) (string, error) {
	if len(priv) == 0 {
		return "", errors.New("sign: private key is empty")
	}
	digest := sha256.Sum256(ciphertext)
	sig := ed25519.Sign(priv, digest[:])
	return hex.EncodeToString(sig), nil
}

// VerifyVault verifies a hex-encoded Ed25519 signature against the SHA-256
// hash of the provided vault ciphertext blob.
func VerifyVault(pub ed25519.PublicKey, ciphertext []byte, sigHex string) error {
	if len(pub) == 0 {
		return errors.New("verify: public key is empty")
	}
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return errors.New("verify: invalid signature encoding")
	}
	digest := sha256.Sum256(ciphertext)
	if !ed25519.Verify(pub, digest[:], sig) {
		return errors.New("verify: signature mismatch — vault may have been tampered with")
	}
	return nil
}
