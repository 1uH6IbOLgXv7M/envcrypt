package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent describes a change detected in a watched .env file.
type WatchEvent struct {
	Path      string
	Checksum  string
	DetectedAt time.Time
}

// WatchEnvFile polls the given .env file at the specified interval and sends
// a WatchEvent on the returned channel whenever the file's checksum changes.
// The caller must close the done channel to stop watching.
func WatchEnvFile(path string, interval time.Duration, done <-chan struct{}) (<-chan WatchEvent, error) {
	last, err := fileChecksum(path)
	if err != nil {
		return nil, fmt.Errorf("watch: initial checksum: %w", err)
	}

	ch := make(chan WatchEvent, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				current, err := fileChecksum(path)
				if err != nil {
					continue
				}
				if current != last {
					last = current
					ch <- WatchEvent{
						Path:       path,
						Checksum:   current,
						DetectedAt: time.Now(),
					}
				}
			}
		}
	}()
	return ch, nil
}

// fileChecksum returns the SHA-256 hex digest of the file at path.
func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
