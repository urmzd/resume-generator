package generators

import (
	"regexp"
	"sort"
	"strings"

	"github.com/urmzd/incipit/resume"
)

// metricPattern matches quantified facts worth bolding. Alternatives are
// ordered most-specific first because Go regexp alternation is leftmost-first:
// dollar amounts ("$2M", "~$0.05", "$1K/month"), percentages including ranges
// ("40%", "45-55%"), multipliers ("4x"), and counts with optional scale
// suffix ("1,500", "2B", "200+"). An approximation prefix (~, <, >) is
// included in the match so "~$2M" bolds as one unit.
var metricPattern = regexp.MustCompile(
	`(?:[~<>]\s?)?\$\d+(?:,\d{3})*(?:\.\d+)?[KMBkmb]?(?:/[A-Za-z]+)?` +
		`|(?:[~<>]\s?)?\b\d+(?:\.\d+)?(?:\s?[-–]\s?\d+(?:\.\d+)?)?%` +
		`|(?:[~<>]\s?)?\b\d+(?:\.\d+)?[x×]\b` +
		`|(?:[~<>]\s?)?\b\d+(?:,\d{3})*(?:\.\d+)?(?:[KMB]\b)?\+?`,
)

// emphasisSpans returns the byte ranges of text to bold, merged and sorted.
// A nil or empty config yields no spans.
func emphasisSpans(text string, emphasis *resume.Emphasis) [][2]int {
	if emphasis == nil {
		return nil
	}

	var spans [][2]int
	if emphasis.Metrics {
		for _, m := range metricPattern.FindAllStringIndex(text, -1) {
			spans = append(spans, [2]int{m[0], m[1]})
		}
	}
	for _, phrase := range emphasis.Phrases {
		phrase = strings.TrimSpace(phrase)
		if phrase == "" {
			continue
		}
		re, err := regexp.Compile(`(?i)` + regexp.QuoteMeta(phrase))
		if err != nil {
			continue
		}
		for _, m := range re.FindAllStringIndex(text, -1) {
			spans = append(spans, [2]int{m[0], m[1]})
		}
	}
	if len(spans) == 0 {
		return nil
	}

	sort.Slice(spans, func(i, j int) bool {
		if spans[i][0] != spans[j][0] {
			return spans[i][0] < spans[j][0]
		}
		return spans[i][1] > spans[j][1]
	})

	merged := spans[:1]
	for _, s := range spans[1:] {
		last := &merged[len(merged)-1]
		if s[0] <= last[1] {
			if s[1] > last[1] {
				last[1] = s[1]
			}
			continue
		}
		merged = append(merged, s)
	}
	return merged
}

// emphasizeText escapes text and wraps emphasized spans in engine markup.
// escape may be nil for engines that render text verbatim (Markdown). With a
// nil or empty emphasis config the output is identical to plain escaping.
func emphasizeText(text string, emphasis *resume.Emphasis, escape func(string) string, open, close string) string {
	if escape == nil {
		escape = func(s string) string { return s }
	}

	spans := emphasisSpans(text, emphasis)
	if len(spans) == 0 {
		return escape(text)
	}

	var b strings.Builder
	prev := 0
	for _, s := range spans {
		b.WriteString(escape(text[prev:s[0]]))
		b.WriteString(open)
		b.WriteString(escape(text[s[0]:s[1]]))
		b.WriteString(close)
		prev = s[1]
	}
	b.WriteString(escape(text[prev:]))
	return b.String()
}
