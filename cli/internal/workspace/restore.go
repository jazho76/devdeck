package workspace

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/ui"
)

const restorePrefix = "__devdeck_restore_"

func Load(p paths.Paths, name string) (Workspace, error) {
	slug := Slugify(name)
	if slug == "" {
		return Workspace{}, fmt.Errorf("workspace name %q has no usable characters", name)
	}
	ws, err := load(p.WorkspaceFile(slug))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Workspace{}, fmt.Errorf("no such workspace: %q", name)
		}
		return Workspace{}, err
	}
	return ws, nil
}

func Restore(p paths.Paths, name string) error {
	ws, err := Load(p, name)
	if err != nil {
		return err
	}
	if ws.SchemaVersion > SchemaVersion {
		return fmt.Errorf("workspace %q was saved by a newer devdeck (schema %d > %d)", name, ws.SchemaVersion, SchemaVersion)
	}
	if len(ws.Sessions) == 0 {
		return fmt.Errorf("workspace %q has no sessions to restore", name)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	attached := os.Getenv("TMUX") != ""
	existing := currentSessions()
	baseIdx := tmuxIndexOption("-gv", "base-index")
	paneBaseIdx := tmuxIndexOption("-gwv", "pane-base-index")

	temps := make([]string, len(ws.Sessions))
	for i := range ws.Sessions {
		temps[i] = fmt.Sprintf("%s%d", restorePrefix, i)
	}

	cleanup := func() {
		for _, t := range temps {
			_ = tmuxRun("kill-session", "-t", t)
		}
	}
	for i, s := range ws.Sessions {
		if err := buildSession(temps[i], s, baseIdx, paneBaseIdx, home); err != nil {
			cleanup()
			return fmt.Errorf("building session %q: %w", s.Name, err)
		}
	}

	focus := 0
	for i, s := range ws.Sessions {
		if s.Attached {
			focus = i
			break
		}
	}

	if attached {
		// Killing the session that hosts this process closes our pane and SIGHUPs us;
		// ignore it so the kill/rename/switch below still finishes.
		signal.Ignore(syscall.SIGHUP)
		if err := tmuxRun("switch-client", "-t", temps[focus]); err != nil {
			cleanup()
			return err
		}
	}
	for _, s := range existing {
		if err := tmuxRun("kill-session", "-t", s); err != nil {
			ui.Warn("could not kill session %q: %v", s, err)
		}
	}
	for i, s := range ws.Sessions {
		if err := tmuxRun("rename-session", "-t", temps[i], s.Name); err != nil {
			return err
		}
	}

	ui.Info("Restored workspace: %s", ws.Name)
	focusName := ws.Sessions[focus].Name
	if attached {
		return tmuxRun("switch-client", "-t", focusName)
	}
	return tmuxStream("attach-session", "-t", focusName)
}

func buildSession(temp string, s Session, baseIdx, paneBaseIdx int, home string) error {
	for wi, w := range s.Windows {
		winTarget := fmt.Sprintf("%s:%d", temp, baseIdx+wi)

		firstCwd := home
		if len(w.Panes) > 0 {
			firstCwd = resolveCwd(w.Panes[0].Cwd, home)
		}

		if wi == 0 {
			args := []string{"new-session", "-d", "-s", temp, "-n", w.Name, "-c", firstCwd}
			if w.Width > 0 && w.Height > 0 {
				args = append(args, "-x", strconv.Itoa(w.Width), "-y", strconv.Itoa(w.Height))
			}
			if err := tmuxRun(args...); err != nil {
				return err
			}
		} else if err := tmuxRun("new-window", "-d", "-t", temp, "-n", w.Name, "-c", firstCwd); err != nil {
			return err
		}

		for _, pane := range w.Panes[1:] {
			if err := tmuxRun("split-window", "-t", winTarget, "-c", resolveCwd(pane.Cwd, home)); err != nil {
				return err
			}
			if err := tmuxRun("select-layout", "-t", winTarget, "tiled"); err != nil {
				return err
			}
		}

		if w.Layout != "" {
			if err := tmuxRun("select-layout", "-t", winTarget, w.Layout); err != nil {
				return err
			}
		}

		for pi, pane := range w.Panes {
			if len(pane.Command) == 0 {
				continue
			}
			target := fmt.Sprintf("%s.%d", winTarget, paneBaseIdx+pi)
			if err := tmuxRun("send-keys", "-t", target, shellQuote(pane.Command), "C-m"); err != nil {
				return err
			}
		}

		for pi, pane := range w.Panes {
			if pane.Active {
				if err := tmuxRun("select-pane", "-t", fmt.Sprintf("%s.%d", winTarget, paneBaseIdx+pi)); err != nil {
					return err
				}
				break
			}
		}
	}

	for wi, w := range s.Windows {
		if w.Active {
			return tmuxRun("select-window", "-t", fmt.Sprintf("%s:%d", temp, baseIdx+wi))
		}
	}
	return nil
}

func shellQuote(args []string) string {
	quoted := make([]string, len(args))
	for i, a := range args {
		quoted[i] = "'" + strings.ReplaceAll(a, "'", `'\''`) + "'"
	}
	return strings.Join(quoted, " ")
}

func resolveCwd(cwd, home string) string {
	if cwd != "" {
		if info, err := os.Stat(cwd); err == nil && info.IsDir() {
			return cwd
		}
	}
	ui.Warn("workspace cwd %q is unavailable, using %s", cwd, home)
	return home
}

func currentSessions() []string {
	out, err := tmuxQuery("list-sessions", "-F", "#{session_name}")
	if err != nil {
		return nil
	}
	out = strings.TrimRight(out, "\n")
	if out == "" {
		return nil
	}
	return strings.Split(out, "\n")
}

func tmuxIndexOption(scope, name string) int {
	out, err := tmuxQuery("show-options", scope, name)
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0
	}
	return n
}
