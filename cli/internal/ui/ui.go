package ui

import (
	"fmt"
	"io"
)

type Logger struct {
	w io.Writer
}

func New(w io.Writer) Logger {
	return Logger{w: w}
}

func (l Logger) Step(label string) {
	fmt.Fprintf(l.w, "\n==> %s\n", label)
}

func (l Logger) Info(format string, a ...any) {
	fmt.Fprintf(l.w, format+"\n", a...)
}

func (l Logger) Warn(format string, a ...any) {
	fmt.Fprintf(l.w, "Warning: "+format+"\n", a...)
}
