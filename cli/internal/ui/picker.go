package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	glyphSelector = "❯"
	glyphFilled   = "◉"
	glyphEmpty    = "○"
)

var pickerInstructions = []string{
	"↑/↓ select",
	"Space: Toggle selection",
	"a: Select all · n: Deselect all",
	"Enter: Confirm · q: Cancel",
}

type Selection struct {
	Selected  map[string]bool
	Cancelled bool
}

func MultiSelect(prompt string, items []string, selected map[string]bool) (Selection, error) {
	final, err := tea.NewProgram(picker{
		prompt:   prompt,
		items:    items,
		selected: cloneBoolSet(selected),
	}).Run()
	if err != nil {
		return Selection{}, err
	}

	m := final.(picker)
	if m.cancelled {
		return Selection{Cancelled: true}, nil
	}
	return Selection{Selected: m.selected}, nil
}

type picker struct {
	prompt    string
	items     []string
	selected  map[string]bool
	cursor    int
	cancelled bool
}

func (m picker) Init() tea.Cmd { return nil }

func (m picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.selected[name] = !m.selected[name]
	case "a":
		for _, name := range m.items {
			m.selected[name] = true
		}
	case "n":
		m.selected = map[string]bool{}
	case "enter":
		return m, tea.Quit
	case "q", "ctrl+c", "esc":
		m.cancelled = true
		return m, tea.Quit
	}
	return m, nil
}

func (m picker) View() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s %s %s\n",
		paint(colorBlue, "?"),
		paint(styleBold+colorWhite, m.prompt),
		paint(colorGrey, "›"),
	)

	for _, line := range pickerInstructions {
		fmt.Fprintf(&b, "  %s\n", paint(colorGrey, line))
	}
	b.WriteString("\n")

	for i, name := range m.items {
		selector := " "
		if i == m.cursor {
			selector = paint(colorBlue, glyphSelector)
		}
		box := glyphEmpty
		if m.selected[name] {
			box = paint(colorGreen, glyphFilled)
		}
		fmt.Fprintf(&b, "%s %s %s\n", selector, box, name)
	}

	return b.String()
}

func cloneBoolSet(s map[string]bool) map[string]bool {
	out := make(map[string]bool, len(s))
	for k, v := range s {
		if v {
			out[k] = true
		}
	}
	return out
}
