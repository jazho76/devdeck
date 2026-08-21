package run

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

func TestForegroundExitCode(t *testing.T) {
	cases := []struct {
		name   string
		script string
		want   int
	}{
		{name: "success", script: "exit 0", want: 0},
		{name: "failure keeps the child status", script: "exit 42", want: 42},
		{name: "signal death reports 128 plus the signal", script: "kill -TERM $$", want: 128 + int(syscall.SIGTERM)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Foreground("sh", "-c", tc.script)
			if err != nil {
				t.Fatalf("Foreground: %v", err)
			}
			if got != tc.want {
				t.Fatalf("exit code = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestForegroundLeavesChildInterruptible(t *testing.T) {
	status := filepath.Join(t.TempDir(), "status")
	if _, err := Foreground("sh", "-c", "grep '^SigIgn:' /proc/self/status > "+status); err != nil {
		t.Fatalf("Foreground: %v", err)
	}

	data, err := os.ReadFile(status)
	if err != nil {
		t.Fatalf("reading child status: %v", err)
	}
	field := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(data)), "SigIgn:"))
	ignored, err := strconv.ParseUint(field, 16, 64)
	if err != nil {
		t.Fatalf("parsing SigIgn %q: %v", field, err)
	}

	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGQUIT} {
		if ignored&(1<<(uint(sig)-1)) != 0 {
			t.Fatalf("child inherited an ignored %v (SigIgn=%016x) and can no longer install its own handler", sig, ignored)
		}
	}
}
