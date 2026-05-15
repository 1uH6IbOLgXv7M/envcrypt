package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"envcrypt/crypto"
)

// runPolicyShow prints the current policy for the working directory.
func runPolicyShow(args []string) error {
	dir, _ := os.Getwd()
	p, err := crypto.LoadPolicy(dir)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}
	fmt.Printf("min_recipients:    %d\n", p.MinRecipients)
	fmt.Printf("require_sign:      %v\n", p.RequireSign)
	fmt.Printf("max_versions:      %d\n", p.MaxVersions)
	fmt.Printf("require_audit_log: %v\n", p.RequireAuditLog)
	if len(p.AllowedKeys) > 0 {
		fmt.Printf("allowed_keys:      %s\n", strings.Join(p.AllowedKeys, ", "))
	} else {
		fmt.Println("allowed_keys:      (any)")
	}
	return nil
}

// runPolicySet updates a single policy field by key=value.
// Usage: envcrypt policy set <field> <value>
func runPolicySet(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: policy set <field> <value>")
	}
	field, value := args[0], args[1]
	dir, _ := os.Getwd()
	p, err := crypto.LoadPolicy(dir)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}
	switch field {
	case "min_recipients":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for min_recipients: %q", value)
		}
		p.MinRecipients = n
	case "require_sign":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid value for require_sign: %q", value)
		}
		p.RequireSign = b
	case "max_versions":
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return fmt.Errorf("invalid value for max_versions: %q", value)
		}
		p.MaxVersions = n
	case "require_audit_log":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid value for require_audit_log: %q", value)
		}
		p.RequireAuditLog = b
	case "allowed_keys":
		if value == "" || value == "any" {
			p.AllowedKeys = nil
		} else {
			p.AllowedKeys = strings.Split(value, ",")
			for i, k := range p.AllowedKeys {
				p.AllowedKeys[i] = strings.TrimSpace(k)
			}
		}
	default:
		return fmt.Errorf("unknown policy field: %q", field)
	}
	if err := crypto.SavePolicy(dir, p); err != nil {
		return fmt.Errorf("save policy: %w", err)
	}
	fmt.Printf("policy: set %s = %s\n", field, value)
	return nil
}
