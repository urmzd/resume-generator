package generators

import (
	"testing"

	"github.com/urmzd/incipit/resume"
)

func TestEmphasizeTextMetrics(t *testing.T) {
	emphasis := &resume.Emphasis{Metrics: true}
	wrap := func(text string) string {
		return emphasizeText(text, emphasis, nil, "<b>", "</b>")
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"dollar with scale and unit", "for under $1K/month, with tracing", "for under <b>$1K/month</b>, with tracing"},
		{"approx dollar", "contributed ~$2M ACV", "contributed <b>~$2M</b> ACV"},
		{"decimal dollar", "at ~$0.05 per answer", "at <b>~$0.05</b> per answer"},
		{"scaled count", "scaled to 2B tokens/month", "scaled to <b>2B</b> tokens/month"},
		{"comma count", "rolled out to 1,500 employees", "rolled out to <b>1,500</b> employees"},
		{"plus count", "a 200+ scenario suite", "a <b>200+</b> scenario suite"},
		{"percent", "about ~40% of the org", "about <b>~40%</b> of the org"},
		{"percent range", "45-55% lower cost", "<b>45-55%</b> lower cost"},
		{"multiplier", "~4x lower latency", "<b>~4x</b> lower latency"},
		{"plain count", "across 10 farms in 5 months", "across <b>10</b> farms in <b>5</b> months"},
		{"digit inside word untouched", "stored in S3 buckets", "stored in S3 buckets"},
		{"no metrics", "Partnered with product and security teams.", "Partnered with product and security teams."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrap(tt.in); got != tt.want {
				t.Errorf("emphasizeText(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestEmphasizeTextPhrases(t *testing.T) {
	emphasis := &resume.Emphasis{Phrases: []string{"MCP", "knowledge graph"}}
	got := emphasizeText("Built the mcp layer over a Knowledge Graph.", emphasis, nil, "<b>", "</b>")
	want := "Built the <b>mcp</b> layer over a <b>Knowledge Graph</b>."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEmphasizeTextOverlappingSpansMerge(t *testing.T) {
	emphasis := &resume.Emphasis{Metrics: true, Phrases: []string{"$2M ACV"}}
	got := emphasizeText("contributed ~$2M ACV overall", emphasis, nil, "<b>", "</b>")
	want := "contributed <b>~$2M ACV</b> overall"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEmphasizeTextNilConfigMatchesEscape(t *testing.T) {
	in := "100% of $5M & more"
	escape := func(s string) string { return s + "!" }
	if got := emphasizeText(in, nil, escape, "<b>", "</b>"); got != escape(in) {
		t.Errorf("nil emphasis: got %q, want plain escape %q", got, escape(in))
	}
	empty := &resume.Emphasis{}
	if got := emphasizeText(in, empty, escape, "<b>", "</b>"); got != escape(in) {
		t.Errorf("empty emphasis: got %q, want plain escape %q", got, escape(in))
	}
}

func TestEmphasizeEscapesInsideSpans(t *testing.T) {
	f := newLaTeXFormatter()
	layout := &resume.Layout{Emphasis: &resume.Emphasis{Metrics: true}}
	got := f.Emphasize("saved ~$2M & 40% effort", layout)
	want := `saved \textbf{\textasciitilde{}\$2M} \& \textbf{40\%} effort`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEmphasizeHTMLEscapes(t *testing.T) {
	f := newHTMLFormatter()
	layout := &resume.Layout{Emphasis: &resume.Emphasis{Phrases: []string{"<script>"}}}
	got := string(f.Emphasize("a <script> tag", layout))
	want := "a <strong>&lt;script&gt;</strong> tag"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEmphasizeMarkdown(t *testing.T) {
	f := newMarkdownFormatter()
	layout := &resume.Layout{Emphasis: &resume.Emphasis{Metrics: true}}
	got := f.Emphasize("shipped 3 services", layout)
	want := "shipped **3** services"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	// Nil layout must leave text verbatim, matching pre-emphasis rendering.
	if got := f.Emphasize("plain *text*", nil); got != "plain *text*" {
		t.Errorf("nil layout: got %q", got)
	}
}
