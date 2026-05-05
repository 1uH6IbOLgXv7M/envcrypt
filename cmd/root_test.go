package cmd

import (
	"os"
	"testing"
)

func TestExecute_NoArgs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"envcrypt"}
	if err := Execute(); err != nil {
		t.Errorf("expected no error for no args, got: %v", err)
	}
}

func TestExecute_UnknownCommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"envcrypt", "foobar"}
	if err := Execute(); err == nil {
		t.Error("expected error for unknown command, got nil")
	}
}

func TestRunVersion(t *testing.T) {
	if err := runVersion(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunHelp(t *testing.T) {
	if err := runHelp(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunEncrypt_NoArgs(t *testing.T) {
	if err := runEncrypt(nil); err == nil {
		t.Error("expected error when no file provided to encrypt")
	}
}

func TestRunEncrypt_WithArg(t *testing.T) {
	if err := runEncrypt([]string{".env"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunDecrypt_NoArgs(t *testing.T) {
	if err := runDecrypt(nil); err == nil {
		t.Error("expected error when no file provided to decrypt")
	}
}

func TestRunDecrypt_WithArg(t *testing.T) {
	if err := runDecrypt([]string{".env.age"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
