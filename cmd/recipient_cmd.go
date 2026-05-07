package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"envcrypt/crypto"
)

const defaultRecipientsFile = ".env-recipients"

// runRecipientAdd adds a named public key to the recipients file.
// Usage: envcrypt recipient add <alias> <public-key>
func runRecipientAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envcrypt recipient add <alias> <public-key>")
	}
	alias, pubKey := args[0], args[1]

	path := recipientsFilePath()
	rf, err := crypto.LoadRecipientsFile(path)
	if err != nil {
		return fmt.Errorf("load recipients: %w", err)
	}
	rf.AddRecipient(alias, pubKey)
	if err := crypto.SaveRecipientsFile(path, rf); err != nil {
		return fmt.Errorf("save recipients: %w", err)
	}
	fmt.Printf("Added recipient %q\n", alias)
	return nil
}

// runRecipientList lists all recipients in the recipients file.
// Usage: envcrypt recipient list
func runRecipientList(args []string) error {
	path := recipientsFilePath()
	rf, err := crypto.LoadRecipientsFile(path)
	if err != nil {
		return fmt.Errorf("load recipients: %w", err)
	}
	if len(rf.Recipients) == 0 {
		fmt.Println("No recipients found.")
		return nil
	}
	maxAlias := 5
	for _, r := range rf.Recipients {
		if len(r.Alias) > maxAlias {
			maxAlias = len(r.Alias)
		}
	}
	fmt.Printf("%-*s  %s\n", maxAlias, "ALIAS", "PUBLIC KEY")
	fmt.Println(strings.Repeat("-", maxAlias+2+20))
	for _, r := range rf.Recipients {
		fmt.Printf("%-*s  %s\n", maxAlias, r.Alias, r.PublicKey)
	}
	return nil
}

// runRecipientRemove removes a recipient by alias.
// Usage: envcrypt recipient remove <alias>
func runRecipientRemove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt recipient remove <alias>")
	}
	alias := args[0]
	path := recipientsFilePath()
	rf, err := crypto.LoadRecipientsFile(path)
	if err != nil {
		return fmt.Errorf("load recipients: %w", err)
	}
	if !rf.RemoveRecipient(alias) {
		return fmt.Errorf("recipient %q not found", alias)
	}
	if err := crypto.SaveRecipientsFile(path, rf); err != nil {
		return fmt.Errorf("save recipients: %w", err)
	}
	fmt.Printf("Removed recipient %q\n", alias)
	return nil
}

func recipientsFilePath() string {
	if p := os.Getenv("ENVCRYPT_RECIPIENTS_FILE"); p != "" {
		return p
	}
	return filepath.Join(".", defaultRecipientsFile)
}
