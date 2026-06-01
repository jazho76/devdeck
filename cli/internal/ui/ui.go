package ui

import (
	"fmt"
	"os"
)

func Step(label string) {
	fmt.Printf("\n%s\n", stepStyle.Render("==> "+label))
}

func Info(format string, a ...any) {
	fmt.Println(infoStyle.Render(fmt.Sprintf(format, a...)))
}

func Warn(format string, a ...any) {
	fmt.Println(warnStyle.Render("Warning: " + fmt.Sprintf(format, a...)))
}

func Error(format string, a ...any) {
	fmt.Fprintln(os.Stderr, errorStyle.Render("Error: "+fmt.Sprintf(format, a...)))
}
