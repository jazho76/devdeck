package toolsets

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Result struct {
	Chosen    map[string]bool
	Cancelled bool
}

func RunPicker(available []string, enabled map[string]bool) (Result, error) {
	final, err := tea.NewProgram(model{items: available, enabled: cloneSet(enabled)}).Run()
	if err != nil {
		return Result{}, err
	}

	m := final.(model)
	if m.cancel {
		return Result{Cancelled: true}, nil
	}
	return Result{Chosen: m.enabled}, nil
}

type model struct {
	items   []string
	enabled map[string]bool
	cursor  int
	cancel  bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	n := len(m.items)
	switch key.String() {
	case "up", "k":
		m.cursor = (m.cursor - 1 + n) % n
	case "down", "j":
		m.cursor = (m.cursor + 1) % n
	case " ":
		name := m.items[m.cursor]
		m.enabled[name] = !m.enabled[name]
	case "a":
		for _, name := range m.items {
			m.enabled[name] = true
		}
	case "n":
		m.enabled = map[string]bool{}
	case "enter":
		return m, tea.Quit
	case "q", "ctrl+c", "esc":
		m.cancel = true
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString("Toolsets  —  ↑/↓ move · space toggle · a all · n none · enter save · q cancel\n\n")
	for i, name := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		mark := " "
		if m.enabled[name] {
			mark = "x"
		}
		fmt.Fprintf(&b, "%s [%s] %s\n", cursor, mark, name)
	}
	return b.String()
}

func cloneSet(s map[string]bool) map[string]bool {
	out := make(map[string]bool, len(s))
	for k, v := range s {
		if v {
			out[k] = true
		}
	}
	return out
}
