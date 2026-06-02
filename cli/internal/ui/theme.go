package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	tokyoBg      = lipgloss.Color("#1a1b26")
	tokyoBgHigh  = lipgloss.Color("#292e42")
	tokyoFg      = lipgloss.Color("#c0caf5")
	tokyoComment = lipgloss.Color("#565f89")
	tokyoBlue    = lipgloss.Color("#7aa2f7")
	tokyoMagenta = lipgloss.Color("#bb9af7")
	tokyoGreen   = lipgloss.Color("#9ece6a")
	tokyoYellow  = lipgloss.Color("#e0af68")
	tokyoRed     = lipgloss.Color("#f7768e")
)

var (
	primaryColor = tokyoBlue
	warnColor    = tokyoYellow
	errorColor   = tokyoRed
)

var (
	stepStyle  = lipgloss.NewStyle().Bold(true).Foreground(primaryColor)
	infoStyle  = lipgloss.NewStyle()
	warnStyle  = lipgloss.NewStyle().Foreground(warnColor)
	errorStyle = lipgloss.NewStyle().Foreground(errorColor)
)

func formTheme() *huh.Theme { return tokyoNightTheme() }

func tokyoNightTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(tokyoBlue).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(tokyoBlue).Bold(true)
	t.Focused.Directory = t.Focused.Directory.Foreground(tokyoBlue)
	t.Focused.Description = t.Focused.Description.Foreground(tokyoYellow)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(tokyoRed)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(tokyoRed)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(tokyoMagenta)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(tokyoMagenta)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(tokyoMagenta)
	t.Focused.Option = t.Focused.Option.Foreground(tokyoFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(tokyoMagenta)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(tokyoGreen)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(tokyoGreen).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(tokyoComment).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(tokyoFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(tokyoBg).Background(tokyoBlue).Bold(true)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(tokyoFg).Background(tokyoBgHigh)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(tokyoMagenta)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(tokyoComment)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(tokyoMagenta)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	t.Help.ShortKey = t.Help.ShortKey.Foreground(tokyoComment)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(tokyoComment)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(tokyoComment)
	t.Help.FullKey = t.Help.FullKey.Foreground(tokyoComment)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(tokyoComment)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(tokyoComment)

	return t
}

func formKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit.SetKeys("esc", "ctrl+c")
	km.Quit.SetHelp("esc", "cancel")
	km.MultiSelect.SelectAll.SetKeys("a")
	km.MultiSelect.SelectAll.SetHelp("a", "select all")
	km.MultiSelect.SelectNone.SetKeys("a")
	km.MultiSelect.SelectNone.SetHelp("a", "select none")
	return km
}
