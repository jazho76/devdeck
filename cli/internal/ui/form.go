package ui

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
)

func sizedForm(field huh.Field) *huh.Form {
	form := huh.NewForm(huh.NewGroup(field)).
		WithTheme(formTheme()).
		WithKeyMap(formKeyMap())
	if w, h, err := term.GetSize(os.Stdout.Fd()); err == nil {
		if w > 0 {
			form = form.WithWidth(w)
		}
		if h > 1 {
			form = form.WithHeight(h - 1)
		}
	}
	return form
}
