package workspace

import (
	"testing"
	"time"

	"github.com/jazho76/devdeck/cli/internal/paths"
)

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
