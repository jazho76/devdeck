package workspace

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const titleBudget = 40

func RestoreLabels(saved []Workspace, now time.Time) []string {
	nameWidth := 0
	for _, w := range saved {
		if n := utf8.RuneCountInString(w.Name); n > nameWidth {
			nameWidth = n
		}
	}

	labels := make([]string, len(saved))
	for i, w := range saved {
		labels[i] = formatRestoreLabel(w, nameWidth, now)
	}
	return labels
}

func formatRestoreLabel(w Workspace, nameWidth int, now time.Time) string {
	session, others := w.activeSession()

	parts := make([]string, 0, 3)
	if summary := windowSummary(session); summary != "" {
		parts = append(parts, summary)
	}
	if others > 0 {
		parts = append(parts, sessionsSuffix(others))
	}
	parts = append(parts, humanizeSince(w.UpdatedAt, now))

	return fmt.Sprintf("%s   %s", padRight(w.Name, nameWidth), strings.Join(parts, " · "))
}

func (w Workspace) activeSession() (Session, int) {
	if len(w.Sessions) == 0 {
		return Session{}, 0
	}
	idx := 0
	for i, s := range w.Sessions {
		if s.Attached {
			idx = i
			break
		}
	}
	return w.Sessions[idx], len(w.Sessions) - 1
}

func windowSummary(s Session) string {
	var b strings.Builder
	for _, w := range s.Windows {
		if w.Name == "" {
			continue
		}
		sep := ""
		if b.Len() > 0 {
			sep = ", "
		}
		if b.Len()+len(sep)+len(w.Name) > titleBudget {
			b.WriteString("…")
			break
		}
		b.WriteString(sep)
		b.WriteString(w.Name)
	}
	return b.String()
}

func sessionsSuffix(others int) string {
	if others == 1 {
		return "+1 session"
	}
	return fmt.Sprintf("+%d sessions", others)
}

func humanizeSince(t, now time.Time) string {
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Local().Format("2006-01-02")
	}
}

func padRight(s string, width int) string {
	if n := utf8.RuneCountInString(s); n < width {
		return s + strings.Repeat(" ", width-n)
	}
	return s
}
