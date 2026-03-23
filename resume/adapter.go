package resume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InputData represents resume data that can be validated and converted to the runtime format.
type InputData interface {
	ToResume() *Resume
	GetFormat() string
	GetSectionOrder() []string
	Validate() error
}

// sectionKeys are the top-level keys that map to renderable sections.
var sectionKeys = map[string]bool{
	"summary":        true,
	"certifications": true,
	"skills":         true,
	"experience":     true,
	"projects":       true,
	"education":      true,
	"languages":      true,
}

// defaultSectionOrder is the fallback when no order can be detected.
var defaultSectionOrder = []string{
	"summary",
	"certifications",
	"experience",
	"education",
	"skills",
	"projects",
	"languages",
}

// ResumeAdapter implements InputData for Resume structures.
type ResumeAdapter struct {
	Resume           *Resume
	SerializationFmt string
	SectionOrder     []string
}

func (a *ResumeAdapter) ToResume() *Resume         { return a.Resume }
func (a *ResumeAdapter) GetFormat() string         { return a.SerializationFmt }
func (a *ResumeAdapter) GetSectionOrder() []string { return a.SectionOrder }

func (a *ResumeAdapter) Validate() error {
	errors := Validate(a.Resume)
	if len(errors) > 0 {
		return fmt.Errorf("validation failed with %d errors: %v", len(errors), errors[0].Message)
	}
	return nil
}

// LoadResumeFromBytes parses resume data from raw bytes with the given format.
// Format must be one of: "json", "md", "markdown", "txt".
func LoadResumeFromBytes(data []byte, format string) (InputData, error) {
	var resumeData Resume
	var serializationFmt string
	var sectionOrder []string

	switch strings.ToLower(format) {
	case "json":
		if err := json.Unmarshal(data, &resumeData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		sectionOrder = extractJSONKeyOrder(data)
		serializationFmt = "json"

	case "md", "markdown", "txt":
		parsed, mdOrder, err := parseMarkdown(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Markdown: %w", err)
		}
		resumeData = *parsed
		sectionOrder = mdOrder
		serializationFmt = "md"

	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: json, md, markdown, txt)", format)
	}

	if len(sectionOrder) == 0 {
		sectionOrder = defaultSectionOrder
	}

	if resumeData.Contact.Name == "" {
		return nil, fmt.Errorf("contact.name is required")
	}

	return &ResumeAdapter{
		Resume:           &resumeData,
		SerializationFmt: serializationFmt,
		SectionOrder:     sectionOrder,
	}, nil
}

// LoadResumeFromFile loads a resume from a JSON or Markdown file.
func LoadResumeFromFile(filePath string) (InputData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	switch strings.ToLower(ext) {
	case "json":
		// keep as-is
	case "md", "markdown", "txt":
		ext = "md"
	default:
		ext = "md"
	}

	return LoadResumeFromBytes(data, ext)
}

// extractJSONKeyOrder uses json.Decoder tokens to find the order of top-level
// keys that correspond to renderable sections.
func extractJSONKeyOrder(data []byte) []string {
	dec := json.NewDecoder(bytes.NewReader(data))

	t, err := dec.Token()
	if err != nil {
		return nil
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return nil
	}

	var order []string
	depth := 0
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			break
		}
		switch v := t.(type) {
		case json.Delim:
			switch v {
			case '{', '[':
				depth++
			case '}', ']':
				depth--
			}
		case string:
			if depth == 0 {
				if sectionKeys[v] {
					order = append(order, v)
				}
				var raw json.RawMessage
				if err := dec.Decode(&raw); err != nil {
					return order
				}
			}
		}
	}
	return order
}
