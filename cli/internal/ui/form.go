package ui

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
)

type Option func(*formConfig)

type formConfig struct {
	fillHeight bool
}

func FillHeight() Option {
	return func(c *formConfig) { c.fillHeight = true }
}

func sizedForm(field huh.Field, opts ...Option) *huh.Form {
	var cfg formConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	form := huh.NewForm(huh.NewGroup(field)).
		WithTheme(formTheme()).
		WithKeyMap(formKeyMap())
	if w, h, err := term.GetSize(os.Stdout.Fd()); err == nil {
		if w > 0 {
			form = form.WithWidth(w)
		}
		if cfg.fillHeight && h > 1 {
			form = form.WithHeight(h - 1)
		}
	}
	return form
}
