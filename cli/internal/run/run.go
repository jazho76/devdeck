package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Stream runs a command with the parent's stdio so the user sees live output
// (git clones, plugin installs, headless nvim).
func Stream(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", line(name, args), err)
	}
	return nil
}

// Output runs a command and returns its trimmed standard output.
func Output(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", line(name, args), err)
	}
	return strings.TrimSpace(string(out)), nil
}

func line(name string, args []string) string {
	if len(args) == 0 {
		return name
	}
	return name + " " + strings.Join(args, " ")
}
