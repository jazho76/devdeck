package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Stream(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", line(bin, args), err)
	}
	return nil
}

func Output(bin string, args ...string) (string, error) {
	out, err := exec.Command(bin, args...).Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", line(bin, args), err)
	}
	return strings.TrimSpace(string(out)), nil
}

func line(bin string, args []string) string {
	if len(args) == 0 {
		return bin
	}
	return bin + " " + strings.Join(args, " ")
}
