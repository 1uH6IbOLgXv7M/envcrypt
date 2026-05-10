package crypto

import (
	"crypto/ed25519"
	"testing"
)

func TestGenerateSigningKeyPair(t *testing.T) {
	kp, err := GenerateSigningKeyPair()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kp.PublicKey) != ed25519.PublicKeySize {
		t.Errorf("unexpected public key size: %d", len(kp.PublicKey))
	}
	if len(kp.PrivateKey) != ed25519.PrivateKeySize {
		t.Errorf("unexpected private key size: %d", len(kp.PrivateKey))
	}
}

func TestGenerateSigningKeyPair_Unique(t *testing.T) {
	kp1, _ := GenerateSigningKeyPair()
	kp2, _ := GenerateSigningKeyPair()
	if string(kp1.PublicKey) == string(kp2.PublicKey) {
		t.Error("expected unique key pairs")
	}
}

func TestSignAndVerifyVault_RoundTrip(t *testing.T) {
	kp, err := GenerateSigningKeyPair()
	if err != nil {
		t.Fatalf("keygen error: %v", err)
	}
	data := []byte("encrypted-vault-payload")
	sig, err := SignVault(kp.PrivateKey, data)
	if err != nil {
		t.Fatalf("sign error: %v", err)
	}
	if err := VerifyVault(kp.PublicKey, data, sig); err != nil {
		t.Errorf("verify error: %v", err)
	}
}

func TestVerifyVault_TamperedData(t *testing.T) {
	kp, _ := GenerateSigningKeyPair()
	data := []byte("original-payload")
	sig, _ := SignVault(kp.PrivateKey, data)
	tampered := []byte("tampered-payload")
	if err := VerifyVault(kp.PublicKey, tampered, sig); err == nil {
		t.Error("expected error for tampered data")
	}
}

func TestVerifyVault_InvalidSigHex(t *testing.T) {
	kp, _ := GenerateSigningKeyPair()
	if err := VerifyVault(kp.PublicKey, []byte("data"), "not-hex!!"); err == nil {
		t.Error("expected error for invalid hex signature")
	}
}

func TestSignVault_EmptyKey(t *testing.T) {
	_, err := SignVault(nil, []byte("data"))
	if err == nil {
		t.Error("expected error for empty private key")
	}
}

func TestVerifyVault_EmptyKey(t *testing.T) {
	err := VerifyVault(nil, []byte("data"), "aabbcc")
	if err == nil {
		t.Error("expected error for empty public key")
	}
}
