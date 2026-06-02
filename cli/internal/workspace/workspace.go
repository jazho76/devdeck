package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jazho76/devdeck/cli/internal/fsx"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/ui"
)

const SchemaVersion = 1

const maxNameLen = 255

type Workspace struct {
	SchemaVersion int       `json:"schemaVersion"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	LastOpenedAt  time.Time `json:"lastOpenedAt"`
	TmuxVersion   string    `json:"tmuxVersion"`
	Sessions      []Session `json:"sessions"`
}

type Session struct {
	Name     string   `json:"name"`
	Attached bool     `json:"attached"`
	Windows  []Window `json:"windows"`
}

type Window struct {
	Index  int    `json:"index"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Layout string `json:"layout"`
	Zoomed bool   `json:"zoomed,omitempty"`
	Panes  []Pane `json:"panes"`
}

type Pane struct {
	Index   int      `json:"index"`
	Active  bool     `json:"active"`
	Cwd     string   `json:"cwd"`
	Command []string `json:"command,omitempty"`
}

var slugStrip = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(name string) string {
	s := slugStrip.ReplaceAllString(strings.ToLower(name), "-")
	return strings.Trim(s, "-")
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("workspace name cannot be empty")
	}
	if len(name) > maxNameLen {
		return fmt.Errorf("workspace name too long (max %d characters)", maxNameLen)
	}
	return nil
}

func Save(p paths.Paths, name string) error {
	if strings.TrimSpace(name) == "" {
		session, err := currentSession()
		if err != nil {
			return fmt.Errorf("could not determine current session; is tmux running? %w", err)
		}
		name = session
	}
	if err := ValidateName(name); err != nil {
		return err
	}
	slug := Slugify(name)
	if slug == "" {
		return fmt.Errorf("workspace name %q has no usable characters for a filename", name)
	}

	version, sessions, err := capture()
	if err != nil {
		return err
	}
	if len(sessions) > 0 {
		sessions[0].Name = slug
	}

	now := time.Now().UTC()
	createdAt := now
	var lastOpenedAt time.Time
	if existing, err := load(p.WorkspaceFile(slug)); err == nil {
		createdAt = existing.CreatedAt
		lastOpenedAt = existing.LastOpenedAt
	} else if !errors.Is(err, os.ErrNotExist) {
		ui.Warn("could not read existing workspace %q, overwriting: %v", slug, err)
	}

	ws := Workspace{
		SchemaVersion: SchemaVersion,
		Name:          name,
		Slug:          slug,
		CreatedAt:     createdAt,
		UpdatedAt:     now,
		LastOpenedAt:  lastOpenedAt,
		TmuxVersion:   version,
		Sessions:      sessions,
	}

	return writeWorkspace(p, ws)
}

func writeWorkspace(p paths.Paths, ws Workspace) error {
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(p.Workspaces, 0o755); err != nil {
		return err
	}
	return fsx.WriteFileAtomic(p.WorkspaceFile(ws.Slug), data, 0o644)
}

func LastActivity(ws Workspace) (when time.Time, viaRestore bool) {
	if ws.LastOpenedAt.After(ws.UpdatedAt) {
		return ws.LastOpenedAt, true
	}
	return ws.UpdatedAt, false
}

func List(p paths.Paths) ([]Workspace, error) {
	matches, err := filepath.Glob(filepath.Join(p.Workspaces, "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)

	out := make([]Workspace, 0, len(matches))
	for _, m := range matches {
		ws, err := load(m)
		if err != nil {
			continue
		}
		out = append(out, ws)
	}
	return out, nil
}

func load(path string) (Workspace, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Workspace{}, err
	}
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return Workspace{}, fmt.Errorf("reading %s: %w", path, err)
	}
	return ws, nil
}
