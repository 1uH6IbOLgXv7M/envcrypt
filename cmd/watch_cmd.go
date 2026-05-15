package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"envcrypt/crypto"
)

// runWatch watches a .env file for changes and prints a notification each time
// the file is modified. It blocks until the user interrupts (Ctrl-C).
//
// Usage: envcrypt watch <file> [--interval <ms>]
func runWatch(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt watch <file>")
	}
	path := args[0]

	interval := 500 * time.Millisecond
	for i := 1; i < len(args)-1; i++ {
		if args[i] == "--interval" {
			var ms int
			if _, err := fmt.Sscanf(args[i+1], "%d", &ms); err == nil && ms > 0 {
				interval = time.Duration(ms) * time.Millisecond
			}
		}
	}

	done := make(chan struct{})
	ch, err := crypto.WatchEnvFile(path, interval, done)
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Watching %s (interval: %s) — press Ctrl-C to stop\n", path, interval)

	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return nil
			}
			fmt.Printf("[%s] change detected in %s (sha256: %s)\n",
				ev.DetectedAt.Format("15:04:05"), ev.Path, ev.Checksum[:12])
		case <-sig:
			close(done)
			fmt.Println("\nStopped watching.")
			return nil
		}
	}
}
