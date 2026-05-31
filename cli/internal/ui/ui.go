package ui

import (
	"fmt"
	"os"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGrey   = "\033[90m"
	colorWhite  = "\033[97m"
	styleBold   = "\033[1m"
	colorReset  = "\033[0m"
)

func Step(label string) {
	fmt.Printf("\n==> %s\n", label)
}

func Info(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Warn(format string, a ...any) {
	fmt.Println(paint(colorYellow, "Warning: "+fmt.Sprintf(format, a...)))
}

func Error(format string, a ...any) {
	fmt.Fprintln(os.Stderr, paint(colorRed, "Error: "+fmt.Sprintf(format, a...)))
}

func paint(color, s string) string {
	if os.Getenv("NO_COLOR") != "" {
		return s
	}
	return color + s + colorReset
}
