package ui

import (
	"errors"

	"github.com/charmbracelet/huh"
)

type Selection struct {
	Selected  map[string]bool
	Cancelled bool
}

func MultiSelect(prompt string, items []string, selected map[string]bool, labels map[string]string, opts ...Option) (Selection, error) {
	options := make([]huh.Option[string], len(items))
	values := make([]string, 0, len(items))
	for i, item := range items {
		label := item
		if l, ok := labels[item]; ok {
			label = l
		}
		options[i] = huh.NewOption(label, item).Selected(selected[item])
		if selected[item] {
			values = append(values, item)
		}
	}

	field := huh.NewMultiSelect[string]().
		Title(prompt).
		Options(options...).
		Filterable(true).
		Value(&values)

	err := sizedForm(field, opts...).Run()
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
