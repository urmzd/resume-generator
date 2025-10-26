package generators

import (
	"testing"
	"unicode/utf8"
)

func FuzzSanitizePhone(f *testing.F) {
	for _, seed := range []string{
		"",
		"+15551234567",
		"(555) 123-4567",
		"001-555-EXAMPLE",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		output := sanitizePhone(input)

		for i, r := range output {
			if r != '+' && (r < '0' || r > '9') {
				t.Fatalf("sanitizePhone(%q) produced invalid rune %q at index %d", input, string(r), i)
			}
		}

		if !utf8.ValidString(output) {
			t.Fatalf("sanitizePhone(%q) produced invalid UTF-8 output", input)
		}
	})
}
