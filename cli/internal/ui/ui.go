package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func Step(label string) {
	fmt.Printf("\n%s\n", stepStyle.Render(":: "+label))
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

func StatusOK(label, detail string) {
	printStatus(okStyle, "✓", label, detail)
}

func StatusWarn(label, detail string) {
	printStatus(warnStyle, "!", label, detail)
}

func StatusFail(label, detail string) {
	printStatus(errorStyle, "✗", label, detail)
}

func printStatus(style lipgloss.Style, glyph, label, detail string) {
	line := style.Render(glyph + " " + label)
	if detail != "" {
		line += "  " + detailStyle.Render(detail)
	}
	fmt.Println(line)
}

func Hint(text string) {
	fmt.Println(hintStyle.Render("    " + text))
}
