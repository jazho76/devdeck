package ui

import (
	"errors"

	"github.com/charmbracelet/huh"
)

type Selection struct {
	Selected  map[string]bool
	Cancelled bool
}

func MultiSelect(prompt string, items []string, selected map[string]bool) (Selection, error) {
	options := make([]huh.Option[string], len(items))
	values := make([]string, 0, len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item, item).Selected(selected[item])
		if selected[item] {
			values = append(values, item)
		}
	}

	field := huh.NewMultiSelect[string]().
		Title(prompt).
		Options(options...).
		Filterable(true).
		Value(&values)

	err := huh.NewForm(huh.NewGroup(field)).
		WithTheme(formTheme()).
		WithKeyMap(formKeyMap()).
		Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return Selection{Cancelled: true}, nil
	}
	if err != nil {
		return Selection{}, err
	}

	result := make(map[string]bool, len(values))
	for _, v := range values {
		result[v] = true
	}
	return Selection{Selected: result}, nil
}
