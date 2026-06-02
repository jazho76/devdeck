package workspace

import (
	"os"
	"testing"
	"time"

	"github.com/jazho76/devdeck/cli/internal/paths"
)

func seed(t *testing.T, p paths.Paths, name string) Workspace {
	t.Helper()
	ws := Workspace{
		SchemaVersion: SchemaVersion,
		Name:          name,
		Slug:          Slugify(name),
		Sessions:      []Session{{Name: "main"}},
	}
	if err := writeWorkspace(p, ws); err != nil {
		t.Fatalf("seed %q: %v", name, err)
	}
	return ws
}

func TestDelete(t *testing.T) {
	t.Run("removes the file", func(t *testing.T) {
		p := paths.Paths{Workspaces: t.TempDir()}
		ws := seed(t, p, "Alpha")
		if err := Delete(p, ws.Name); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if _, err := os.Stat(p.WorkspaceFile(ws.Slug)); !os.IsNotExist(err) {
			t.Fatalf("file still present, stat err = %v", err)
		}
	})

	t.Run("missing workspace errors", func(t *testing.T) {
		p := paths.Paths{Workspaces: t.TempDir()}
		if err := Delete(p, "ghost"); err == nil {
			t.Fatal("expected error deleting missing workspace")
		}
	})
}

func TestRename(t *testing.T) {
	t.Run("slug change moves the file", func(t *testing.T) {
		p := paths.Paths{Workspaces: t.TempDir()}
		seed(t, p, "Alpha")
		ws, err := Rename(p, "Alpha", "Beta")
		if err != nil {
			t.Fatalf("Rename: %v", err)
		}
		if ws.Name != "Beta" || ws.Slug != "beta" {
			t.Fatalf("got name=%q slug=%q, want Beta/beta", ws.Name, ws.Slug)
		}
		if _, err := os.Stat(p.WorkspaceFile("alpha")); !os.IsNotExist(err) {
			t.Fatalf("old file still present, stat err = %v", err)
		}
		got, err := load(p.WorkspaceFile("beta"))
		if err != nil {
			t.Fatalf("load beta: %v", err)
		}
		if got.Name != "Beta" {
			t.Fatalf("persisted name = %q, want Beta", got.Name)
		}
	})

	t.Run("casing-only rename updates in place", func(t *testing.T) {
		p := paths.Paths{Workspaces: t.TempDir()}
		seed(t, p, "alpha")
		ws, err := Rename(p, "alpha", "Alpha")
		if err != nil {
			t.Fatalf("Rename: %v", err)
		}
		if ws.Name != "Alpha" || ws.Slug != "alpha" {
			t.Fatalf("got name=%q slug=%q, want Alpha/alpha", ws.Name, ws.Slug)
		}
		got, err := load(p.WorkspaceFile("alpha"))
		if err != nil {
			t.Fatalf("load alpha: %v", err)
		}
		if got.Name != "Alpha" {
			t.Fatalf("persisted name = %q, want Alpha", got.Name)
		}
	})

	t.Run("collision errors and leaves both intact", func(t *testing.T) {
		p := paths.Paths{Workspaces: t.TempDir()}
		seed(t, p, "Alpha")
		seed(t, p, "Beta")
		if _, err := Rename(p, "Alpha", "Beta"); err == nil {
			t.Fatal("expected collision error")
		}
		if _, err := os.Stat(p.WorkspaceFile("alpha")); err != nil {
			t.Fatalf("alpha should survive: %v", err)
		}
		beta, err := load(p.WorkspaceFile("beta"))
		if err != nil {
			t.Fatalf("load beta: %v", err)
		}
		if beta.Name != "Beta" {
			t.Fatalf("beta clobbered: name = %q", beta.Name)
		}
	})
}

func TestWriteWorkspaceRoundTripPreservesLastOpenedAt(t *testing.T) {
	p := paths.Paths{Workspaces: t.TempDir()}
	want := Workspace{
		SchemaVersion: SchemaVersion,
		Name:          "Alpha",
		Slug:          "alpha",
		CreatedAt:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		LastOpenedAt:  time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
	}

	if err := writeWorkspace(p, want); err != nil {
		t.Fatalf("writeWorkspace: %v", err)
	}
	got, err := load(p.WorkspaceFile(want.Slug))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !got.LastOpenedAt.Equal(want.LastOpenedAt) {
		t.Fatalf("LastOpenedAt = %v, want %v", got.LastOpenedAt, want.LastOpenedAt)
	}
}

func TestLastActivity(t *testing.T) {
	updated := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	base := Workspace{Slug: "alpha", UpdatedAt: updated}

	t.Run("zero LastOpenedAt falls back to UpdatedAt", func(t *testing.T) {
		when, viaRestore := LastActivity(base)
		if viaRestore || !when.Equal(updated) {
			t.Fatalf("got (%v, %v), want (%v, false)", when, viaRestore, updated)
		}
	})

	t.Run("older LastOpenedAt falls back to UpdatedAt", func(t *testing.T) {
		ws := base
		ws.LastOpenedAt = updated.Add(-24 * time.Hour)
		when, viaRestore := LastActivity(ws)
		if viaRestore || !when.Equal(updated) {
			t.Fatalf("got (%v, %v), want (%v, false)", when, viaRestore, updated)
		}
	})

	t.Run("newer LastOpenedAt wins", func(t *testing.T) {
		ws := base
		ws.LastOpenedAt = updated.Add(24 * time.Hour)
		when, viaRestore := LastActivity(ws)
		if !viaRestore || !when.Equal(ws.LastOpenedAt) {
			t.Fatalf("got (%v, %v), want (%v, true)", when, viaRestore, ws.LastOpenedAt)
		}
	})
}
