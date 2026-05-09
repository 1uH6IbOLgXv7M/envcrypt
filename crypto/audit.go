package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEvent represents a single recorded action in the audit log.
type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	Details   string    `json:"details"`
}

// AuditLog holds a list of audit events.
type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

// AuditLogPath returns the path to the audit log file within dir.
func AuditLogPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_audit.json")
}

// LoadAuditLog reads the audit log from dir. Returns an empty log if missing.
func LoadAuditLog(dir string) (*AuditLog, error) {
	path := AuditLogPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read audit log: %w", err)
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parse audit log: %w", err)
	}
	return &log, nil
}

// SaveAuditLog writes the audit log to dir.
func SaveAuditLog(dir string, log *AuditLog) error {
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal audit log: %w", err)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create audit dir: %w", err)
	}
	return os.WriteFile(AuditLogPath(dir), data, 0600)
}

// AppendEvent appends a new event to the audit log and saves it.
func AppendEvent(dir, action, actor, details string) error {
	log, err := LoadAuditLog(dir)
	if err != nil {
		return err
	}
	log.Events = append(log.Events, AuditEvent{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Actor:     actor,
		Details:   details,
	})
	return SaveAuditLog(dir, log)
}
