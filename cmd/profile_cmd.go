package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"envcrypt/crypto"
)

func runProfileAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envcrypt profile add <name> <env-file>")
	}
	name, envFile := args[0], args[1]
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	store, err := crypto.LoadProfileStore(dir)
	if err != nil {
		return fmt.Errorf("load profile store: %w", err)
	}
	store.AddProfile(crypto.Profile{
		Name:    name,
		EnvFile: envFile,
	})
	if err := crypto.SaveProfileStore(dir, store); err != nil {
		return fmt.Errorf("save profile store: %w", err)
	}
	fmt.Printf("profile '%s' -> %s added\n", name, envFile)
	return nil
}

func runProfileList(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	store, err := crypto.LoadProfileStore(dir)
	if err != nil {
		return fmt.Errorf("load profile store: %w", err)
	}
	if len(store.Profiles) == 0 {
		fmt.Println("no profiles defined")
		return nil
	}
	names := make([]string, 0, len(store.Profiles))
	for n := range store.Profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tENV FILE\tUPDATED")
	for _, n := range names {
		p := store.Profiles[n]
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.EnvFile, p.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}

func runProfileRemove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt profile remove <name>")
	}
	name := args[0]
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	store, err := crypto.LoadProfileStore(dir)
	if err != nil {
		return fmt.Errorf("load profile store: %w", err)
	}
	if !store.RemoveProfile(name) {
		return fmt.Errorf("profile '%s' not found", name)
	}
	if err := crypto.SaveProfileStore(dir, store); err != nil {
		return fmt.Errorf("save profile store: %w", err)
	}
	fmt.Printf("profile '%s' removed\n", name)
	return nil
}
