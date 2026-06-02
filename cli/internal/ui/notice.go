package ui

import (
	"errors"

	"github.com/charmbracelet/huh"
)

func Notice(title, message string) error {
	field := huh.NewNote().
		Title(title).
		Description(message).
		Next(true).
		NextLabel("OK")

	err := sizedForm(field).Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return nil
	}
	return err
}
