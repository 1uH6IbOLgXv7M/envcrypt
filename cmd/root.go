package cmd

import (
	"errors"
	"fmt"
	"os"
)

const version = "0.1.0"

var validCommands = map[string]func([]string) error{
	"encrypt": runEncrypt,
	"decrypt": runDecrypt,
	"version": runVersion,
	"help":    runHelp,
}

func Execute() error {
	args := os.Args[1:]
	if len(args) == 0 {
		return runHelp(nil)
	}

	command := args[0]
	fn, ok := validCommands[command]
	if !ok {
		return fmt.Errorf("unknown command %q. Run 'envcrypt help' for usage", command)
	}

	return fn(args[1:])
}

func runVersion(_ []string) error {
	fmt.Printf("envcrypt version %s\n", version)
	return nil
}

func runHelp(_ []string) error {
	fmt.Println(`envcrypt - encrypt and version-control .env files using age encryption

Usage:
  envcrypt <command> [arguments]

Commands:
  encrypt   Encrypt a .env file
  decrypt   Decrypt an encrypted .env file
  version   Print the version
  help      Show this help message

Examples:
  envcrypt encrypt .env
  envcrypt decrypt .env.age`)
	return nil
}

func runEncrypt(args []string) error {
	if len(args) < 1 {
		return errors.New("encrypt requires a file argument: envcrypt encrypt <file>")
	}
	fmt.Printf("encrypting %s...\n", args[0])
	return nil
}

func runDecrypt(args []string) error {
	if len(args) < 1 {
		return errors.New("decrypt requires a file argument: envcrypt decrypt <file>")
	}
	fmt.Printf("decrypting %s...\n", args[0])
	return nil
}
