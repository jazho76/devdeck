package workspace

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func capture() (version string, sessions []Session, err error) {
	const sep = "\t"
	format := strings.Join([]string{
		"#{version}",
		"#{session_name}",
		"#{?session_attached,1,0}",
		"#{window_index}",
		"#{window_name}",
		"#{?window_active,1,0}",
		"#{window_width}",
		"#{window_height}",
		"#{window_layout}",
		"#{pane_index}",
		"#{?pane_active,1,0}",
		"#{pane_current_path}",
		"#{pane_current_command}",
		"#{pane_pid}",
	}, sep)

	out, err := tmuxQuery("list-panes", "-a", "-F", format)
	if err != nil {
		return "", nil, fmt.Errorf("could not read tmux server; is tmux running? %w", err)
	}
	out = strings.TrimRight(out, "\n")
	if out == "" {
		return "", nil, nil
	}

	sessionAt := map[string]int{}
	windowAt := map[string]int{}

	for _, raw := range strings.Split(out, "\n") {
		f := strings.Split(raw, sep)
		if len(f) != 14 {
			return "", nil, fmt.Errorf("unexpected tmux output: %q", raw)
		}
		wi, err := strconv.Atoi(f[3])
		if err != nil {
			return "", nil, fmt.Errorf("parsing window index %q: %w", f[3], err)
		}
		ww, err := strconv.Atoi(f[6])
		if err != nil {
			return "", nil, fmt.Errorf("parsing window width %q: %w", f[6], err)
		}
		wh, err := strconv.Atoi(f[7])
		if err != nil {
			return "", nil, fmt.Errorf("parsing window height %q: %w", f[7], err)
		}
		pi, err := strconv.Atoi(f[9])
		if err != nil {
			return "", nil, fmt.Errorf("parsing pane index %q: %w", f[9], err)
		}

		version = f[0]
		sName := f[1]
		si, ok := sessionAt[sName]
		if !ok {
			sessions = append(sessions, Session{Name: sName, Attached: f[2] == "1"})
			si = len(sessions) - 1
			sessionAt[sName] = si
		}

		wKey := sName + "\x00" + f[3]
		wj, ok := windowAt[wKey]
		if !ok {
			sessions[si].Windows = append(sessions[si].Windows, Window{
				Index:  wi,
				Name:   f[4],
				Active: f[5] == "1",
				Width:  ww,
				Height: wh,
				Layout: f[8],
			})
			wj = len(sessions[si].Windows) - 1
			windowAt[wKey] = wj
		}

		var command []string
		if f[12] == nvimCommand {
			if shellPID, err := strconv.Atoi(f[13]); err == nil {
				command = nvimArgs(shellPID)
			}
		}

		sessions[si].Windows[wj].Panes = append(sessions[si].Windows[wj].Panes, Pane{
			Index:   pi,
			Active:  f[10] == "1",
			Cwd:     f[11],
			Command: command,
		})
	}

	return version, sessions, nil
}

const nvimCommand = "nvim"

func nvimArgs(shellPID int) []string {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		name, ppid, ok := procStatus(pid)
		if !ok || name != nvimCommand || ppid != shellPID {
			continue
		}
		return procArgs(pid)
	}
	return nil
}

func procStatus(pid int) (name string, ppid int, ok bool) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return "", 0, false
	}
	var haveName, havePPID bool
	for _, line := range strings.Split(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "Name:"):
			name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			haveName = true
		case strings.HasPrefix(line, "PPid:"):
			ppid, err = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "PPid:")))
			havePPID = err == nil
		}
		if haveName && havePPID {
			return name, ppid, true
		}
	}
	return "", 0, false
}

func procArgs(pid int) []string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil || len(data) == 0 {
		return nil
	}
	const argSeparator = "\x00"
	args := strings.Split(strings.TrimRight(string(data), argSeparator), argSeparator)
	if len(args) == 0 {
		return nil
	}
	return args
}
