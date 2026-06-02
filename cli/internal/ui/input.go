package ui

import (
	"errors"

	"github.com/charmbracelet/huh"
)

type Input struct {
	Value     string
	Cancelled bool
}

func Prompt(prompt, initial string, validate func(string) (hint string, ok bool), opts ...Option) (Input, error) {
	value := initial

	field := huh.NewInput().
		Title(prompt).
		Value(&value).
		Validate(func(s string) error {
			if hint, ok := validate(s); !ok {
				return errors.New(hint)
			}
			return nil
		}).
		DescriptionFunc(func() string {
			if hint, ok := validate(value); ok {
				return hint
			}
			return ""
		}, &value)

	err := sizedForm(field, opts...).Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return Input{Cancelled: true}, nil
	}
	if err != nil {
		return Input{}, err
	}
	return Input{Value: value}, nil
}
