package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"envcrypt/crypto"
)

// runTagAdd adds a named tag pointing to a vault version.
// Usage: envcrypt tag add <name> <version> [message] --dir <dir>
func runTagAdd(args []string, dir string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: tag add <name> <version> [message]")
	}
	name := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}
	message := ""
	if len(args) >= 3 {
		message = args[2]
	}

	// Validate the version exists in the vault.
	vault, err := crypto.LoadVault(dir)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}
	if _, err := vault.EntryByVersion(version); err != nil {
		return fmt.Errorf("version %d not found in vault: %w", version, err)
	}

	store, err := crypto.LoadTagStore(dir)
	if err != nil {
		return fmt.Errorf("load tag store: %w", err)
	}
	store.AddTag(name, version, message)
	if err := crypto.SaveTagStore(dir, store); err != nil {
		return fmt.Errorf("save tag store: %w", err)
	}
	fmt.Printf("Tagged version %d as %q\n", version, name)
	return nil
}

// runTagList lists all tags in the store.
func runTagList(dir string) error {
	store, err := crypto.LoadTagStore(dir)
	if err != nil {
		return fmt.Errorf("load tag store: %w", err)
	}
	if len(store.Tags) == 0 {
		fmt.Println("No tags found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVERSION\tCREATED\tMESSAGE")
	for _, t := range store.Tags {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", t.Name, t.Version, t.CreatedAt.Format("2006-01-02 15:04:05"), t.Message)
	}
	return w.Flush()
}

// runTagRemove removes a named tag.
func runTagRemove(args []string, dir string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: tag remove <name>")
	}
	store, err := crypto.LoadTagStore(dir)
	if err != nil {
		return fmt.Errorf("load tag store: %w", err)
	}
	if !store.RemoveTag(args[0]) {
		return fmt.Errorf("tag %q not found", args[0])
	}
	if err := crypto.SaveTagStore(dir, store); err != nil {
		return fmt.Errorf("save tag store: %w", err)
	}
	fmt.Printf("Removed tag %q\n", args[0])
	return nil
}
