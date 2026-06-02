package ui

import (
	"errors"

	"github.com/charmbracelet/huh"
)

type Choice struct {
	Index     int
	Cancelled bool
}

func SingleSelect(prompt string, items []string, opts ...Option) (Choice, error) {
	options := make([]huh.Option[int], len(items))
	for i, label := range items {
		options[i] = huh.NewOption(label, i)
	}

	var index int
	field := huh.NewSelect[int]().
		Title(prompt).
		Options(options...).
		Value(&index)

	err := sizedForm(field, opts...).Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return Choice{Cancelled: true}, nil
	}
	if err != nil {
		return Choice{}, err
	}
	return Choice{Index: index}, nil
}
