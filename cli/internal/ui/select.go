package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var selectInstructions = []string{
	"↑/↓ select",
	"Enter: Confirm · q: Cancel",
}

type Choice struct {
	Index     int
	Cancelled bool
}

func SingleSelect(prompt string, items []string) (Choice, error) {
	final, err := tea.NewProgram(selector{prompt: prompt, items: items}).Run()
	if err != nil {
		return Choice{}, err
	}
	m := final.(selector)
	if m.cancelled {
		return Choice{Cancelled: true}, nil
	}
	return Choice{Index: m.cursor}, nil
}

type selector struct {
	prompt    string
	items     []string
	cursor    int
	cancelled bool
}

func (m selector) Init() tea.Cmd { return nil }

func (m selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case "enter":
		return m, tea.Quit
	case "q", "ctrl+c", "esc":
		m.cancelled = true
		return m, tea.Quit
	}
	return m, nil
}

func (m selector) View() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s %s\n",
		paint(colorBlue, "?"),
		paint(styleBold+colorWhite, m.prompt),
		paint(colorGrey, "›"),
	)
	for _, line := range selectInstructions {
		fmt.Fprintf(&b, "  %s\n", paint(colorGrey, line))
	}
	b.WriteString("\n")

	for i, name := range m.items {
		if i == m.cursor {
			fmt.Fprintf(&b, "%s %s\n", paint(colorBlue, glyphSelector), paint(colorWhite, name))
			continue
		}
		fmt.Fprintf(&b, "  %s\n", name)
	}
	return b.String()
}
