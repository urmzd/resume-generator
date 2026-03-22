// Package ui provides styled terminal output helpers for the incipit CLI.
// All output is written to stderr so that stdout remains clean for data (e.g. JSON schema).
package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	cyanBold   = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	greenBold  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	yellowBold = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	redBold    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	dim        = lipgloss.NewStyle().Faint(true)
	cyanPlain  = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)

// Header prints a cyan bold command name followed by a 40-char dim horizontal rule.
func Header(cmdName string) {
	fmt.Fprintf(os.Stderr, "  %s\n", cyanBold.Render(cmdName))
	fmt.Fprintf(os.Stderr, "  %s\n", dim.Render(strings.Repeat("─", 40)))
}

// PhaseOk prints a green checkmark with a message and optional dim detail.
func PhaseOk(msg string, detail string) {
	if detail != "" {
		fmt.Fprintf(os.Stderr, "  %s %s %s\n", greenBold.Render("✓"), msg, dim.Render(detail))
	} else {
		fmt.Fprintf(os.Stderr, "  %s %s\n", greenBold.Render("✓"), msg)
	}
}

// Warn prints a yellow warning symbol with a message.
func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "  %s %s\n", yellowBold.Render("⚠"), msg)
}

// Info prints a cyan info symbol with dim text.
func Info(msg string) {
	fmt.Fprintf(os.Stderr, "  %s %s\n", cyanPlain.Render("ℹ"), dim.Render(msg))
}

// Error prints a red error symbol with a message.
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "  %s %s\n", redBold.Render("✗"), msg)
}

// Errorf prints a formatted red error.
func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

// Infof prints a formatted info message.
func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

// Warnf prints a formatted warning.
func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

// Section prints a styled section label (cyan bold, no rule).
func Section(label string) {
	fmt.Fprintf(os.Stderr, "  %s\n", cyanBold.Render(label))
}

// Detail prints a 2-space indented dim detail line.
func Detail(msg string) {
	fmt.Fprintf(os.Stderr, "    %s\n", dim.Render(msg))
}

// Blank prints an empty line to stderr.
func Blank() {
	fmt.Fprintln(os.Stderr)
}
