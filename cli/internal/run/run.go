package run

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func Foreground(bin string, args ...string) (int, error) {
	surviveTerminalInterrupts()

	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return 0, nil
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		return 0, fmt.Errorf("%s: %w", line(bin, args), err)
	}
	if status, ok := exit.Sys().(syscall.WaitStatus); ok && status.Signaled() {
		return 128 + int(status.Signal()), nil
	}
	return exit.ExitCode(), nil
}

func surviveTerminalInterrupts() {
	discarded := make(chan os.Signal, 1)
	signal.Notify(discarded, syscall.SIGINT, syscall.SIGQUIT)
}

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

func Query(bin string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(bin, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("%s: %s", bin, msg)
		}
		return "", fmt.Errorf("%s: %w", bin, err)
	}
	return stdout.String(), nil
}

func line(bin string, args []string) string {
	if len(args) == 0 {
		return bin
	}
	return bin + " " + strings.Join(args, " ")
}
