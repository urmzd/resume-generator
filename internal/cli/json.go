package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/term"
)

// outputJSON marshals v as indented JSON to stdout.
func outputJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		outputError(fmt.Errorf("failed to marshal JSON: %w", err))
		return
	}
	fmt.Println(string(data))
}

// outputError writes a JSON error object to stdout and exits with code 1.
func outputError(err error) {
	data, _ := json.Marshal(map[string]string{"error": err.Error()})
	fmt.Println(string(data))
	os.Exit(1)
}

// stderrLog writes a message to stderr when connected to a TTY.
func stderrLog(msg string) {
	if term.IsTerminal(int(os.Stderr.Fd())) {
		fmt.Fprintln(os.Stderr, msg)
	}
}

