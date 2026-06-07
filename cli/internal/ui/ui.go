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

func StatusHeading(label string) {
	fmt.Println(headingStyle.Render(label))
}

func Dim(s string) string {
	return detailStyle.Render(s)
}

func StatusOKSub(label, detail string)   { printStatusIndented("    ", okStyle, "✓", label, detail) }
func StatusWarnSub(label, detail string) { printStatusIndented("    ", warnStyle, "!", label, detail) }
func StatusFailSub(label, detail string) {
	printStatusIndented("    ", errorStyle, "✗", label, detail)
}

func printStatus(style lipgloss.Style, glyph, label, detail string) {
	printStatusIndented("", style, glyph, label, detail)
}

func printStatusIndented(indent string, style lipgloss.Style, glyph, label, detail string) {
	line := indent + style.Render(glyph+" "+label)
	if detail != "" {
		line += "  " + detailStyle.Render(detail)
	}
	fmt.Println(line)
}

func Hint(text string) {
	fmt.Println(hintStyle.Render("    " + text))
}
