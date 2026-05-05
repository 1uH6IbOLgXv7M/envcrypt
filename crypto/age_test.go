package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}
	if kp.PublicKey == "" || kp.PrivateKey == "" {
		t.Fatal("expected non-empty keys")
	}
	// Keys should be valid base64
	if _, err := base64.StdEncoding.DecodeString(kp.PublicKey); err != nil {
		t.Errorf("PublicKey is not valid base64: %v", err)
	}
	if _, err := base64.StdEncoding.DecodeString(kp.PrivateKey); err != nil {
		t.Errorf("PrivateKey is not valid base64: %v", err)
	}
}

func TestGenerateKeyPair_Unique(t *testing.T) {
	kp1, _ := GenerateKeyPair()
	kp2, _ := GenerateKeyPair()
	if kp1.PrivateKey == kp2.PrivateKey {
		t.Error("two generated key pairs should not be identical")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}

	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123")

	ciphertext, err := Encrypt(plaintext, kp.PublicKey)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if bytes.Equal(ciphertext, plaintext) {
		t.Error("ciphertext should differ from plaintext")
	}

	decrypted, err := Decrypt(ciphertext, kp.PrivateKey)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestEncrypt_InvalidKey(t *testing.T) {
	_, err := Encrypt([]byte("data"), "not-base64!!!")
	if err == nil {
		t.Error("expected error for invalid public key")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	kp, _ := GenerateKeyPair()
	_, err := Decrypt([]byte("short"), kp.PrivateKey)
	if err == nil {
		t.Error("expected error for ciphertext that is too short")
	}
}

func TestDecrypt_InvalidKey(t *testing.T) {
	_, err := Decrypt(make([]byte, 32), "not-base64!!!")
	if err == nil {
		t.Error("expected error for invalid private key")
	}
}
