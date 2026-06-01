package workspace

import (
	"os"

	"github.com/jazho76/devdeck/cli/internal/run"
)

// -u forces tmux to treat its I/O as UTF-8; without it a tmux server running
// under a C/POSIX locale rewrites the tab field separator (and any non-ASCII
// path) to "_", corrupting the capture.
func tmuxQuery(args ...string) (string, error) {
	return run.Query("tmux", append([]string{"-u"}, args...)...)
}

func tmuxRun(args ...string) error {
	_, err := tmuxQuery(args...)
	return err
}

func tmuxStream(args ...string) error {
	return run.Stream("tmux", append([]string{"-u"}, args...)...)
}

func Notify(msg string) {
	if os.Getenv("TMUX") == "" {
		return
	}
	_ = tmuxRun("display-message", msg)
}
