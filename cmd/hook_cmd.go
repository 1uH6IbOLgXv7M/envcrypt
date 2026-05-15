package cmd

import (
	"fmt"
	"os"
	"strings"

	"envcrypt/crypto"
)

// runHookAdd adds or updates a lifecycle hook.
// Usage: envcrypt hook add <event> <command>
func runHookAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envcrypt hook add <event> <command>")
	}
	event := crypto.HookEvent(args[0])
	command := strings.Join(args[1:], " ")

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	hs, err := crypto.LoadHookStore(dir)
	if err != nil {
		return fmt.Errorf("load hook store: %w", err)
	}

	hs.AddHook(event, command)

	if err := crypto.SaveHookStore(dir, hs); err != nil {
		return fmt.Errorf("save hook store: %w", err)
	}

	fmt.Printf("hook registered: [%s] -> %s\n", event, command)
	return nil
}

// runHookList lists all registered hooks.
// Usage: envcrypt hook list
func runHookList(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	hs, err := crypto.LoadHookStore(dir)
	if err != nil {
		return fmt.Errorf("load hook store: %w", err)
	}

	if len(hs.Hooks) == 0 {
		fmt.Println("no hooks registered")
		return nil
	}

	fmt.Printf("%-15s  %s\n", "EVENT", "COMMAND")
	fmt.Println(strings.Repeat("-", 50))
	for _, h := range hs.Hooks {
		status := ""
		if !h.Enabled {
			status = " (disabled)"
		}
		fmt.Printf("%-15s  %s%s\n", h.Event, h.Command, status)
	}
	return nil
}

// runHookRemove removes a lifecycle hook by event.
// Usage: envcrypt hook remove <event>
func runHookRemove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt hook remove <event>")
	}
	event := crypto.HookEvent(args[0])

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	hs, err := crypto.LoadHookStore(dir)
	if err != nil {
		return fmt.Errorf("load hook store: %w", err)
	}

	if !hs.RemoveHook(event) {
		return fmt.Errorf("hook not found for event: %s", event)
	}

	if err := crypto.SaveHookStore(dir, hs); err != nil {
		return fmt.Errorf("save hook store: %w", err)
	}

	fmt.Printf("hook removed: %s\n", event)
	return nil
}
