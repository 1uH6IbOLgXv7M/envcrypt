package crypto

import (
	"testing"
)

func TestMultiRecipientEncryptDecrypt_RoundTrip(t *testing.T) {
	pub1, priv1, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair 1: %v", err)
	}
	pub2, priv2, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair 2: %v", err)
	}

	plaintext := []byte("SECRET=hello\nTOKEN=world")
	recipients := map[string]string{
		"alice": pub1,
		"bob":   pub2,
	}

	ciphertexts, err := MultiRecipientEncrypt(plaintext, recipients)
	if err != nil {
		t.Fatalf("MultiRecipientEncrypt: %v", err)
	}
	if len(ciphertexts) != 2 {
		t.Fatalf("expected 2 ciphertexts, got %d", len(ciphertexts))
	}

	for _, priv := range []string{priv1, priv2} {
		got, alias, err := MultiRecipientDecrypt(ciphertexts, priv)
		if err != nil {
			t.Fatalf("MultiRecipientDecrypt: %v", err)
		}
		if string(got) != string(plaintext) {
			t.Errorf("alias %q: got %q, want %q", alias, got, plaintext)
		}
	}
}

func TestMultiRecipientEncrypt_NoRecipients(t *testing.T) {
	_, err := MultiRecipientEncrypt([]byte("data"), map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty recipients")
	}
}

func TestMultiRecipientDecrypt_NoMatch(t *testing.T) {
	pub1, _, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}
	_, priv2, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair 2: %v", err)
	}

	ciphertexts, err := MultiRecipientEncrypt([]byte("data"), map[string]string{"alice": pub1})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	_, _, err = MultiRecipientDecrypt(ciphertexts, priv2)
	if err == nil {
		t.Fatal("expected error when no matching recipient")
	}
}

func TestMultiRecipientDecrypt_EmptyCiphertexts(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}
	_, _, err = MultiRecipientDecrypt(map[string][]byte{}, priv)
	if err == nil {
		t.Fatal("expected error for empty ciphertexts")
	}
}
