package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	Value     string
	Cancelled bool
}

func Prompt(prompt string, validate func(string) (hint string, ok bool)) (Input, error) {
	final, err := tea.NewProgram(input{prompt: prompt, validate: validate}).Run()
	if err != nil {
		return Input{}, err
	}
	m := final.(input)
	if m.cancelled {
		return Input{Cancelled: true}, nil
	}
	return Input{Value: m.value}, nil
}

type input struct {
	prompt    string
	value     string
	validate  func(string) (string, bool)
	cancelled bool
}

func (m input) Init() tea.Cmd { return nil }

func (m input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.Type {
	case tea.KeyRunes:
		m.value += string(key.Runes)
	case tea.KeySpace:
		m.value += " "
	case tea.KeyBackspace:
		if m.value != "" {
			r := []rune(m.value)
			m.value = string(r[:len(r)-1])
		}
	case tea.KeyEnter:
		if _, valid := m.validate(m.value); valid {
			return m, tea.Quit
		}
	case tea.KeyEsc, tea.KeyCtrlC:
		m.cancelled = true
		return m, tea.Quit
	}
	return m, nil
}

func (m input) View() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s %s %s\n",
		paint(colorBlue, "?"),
		paint(styleBold+colorWhite, m.prompt),
		paint(colorGrey, "›"),
	)
	fmt.Fprintf(&b, "  %s\n\n", paint(colorGrey, "Enter: save · Esc: cancel"))

	fmt.Fprintf(&b, "%s %s%s\n", paint(colorBlue, glyphSelector), m.value, paint(colorGrey, "▮"))

	if m.value != "" {
		if hint, valid := m.validate(m.value); hint != "" {
			color := colorYellow
			if !valid {
				color = colorRed
			}
			fmt.Fprintf(&b, "%s\n", paint(color, hint))
		}
	}

	return b.String()
}
