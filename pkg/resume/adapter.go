package resume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// InputData represents resume data that can be validated and converted to the runtime format.
type InputData interface {
	// ToResume returns the Resume representation
	ToResume() *Resume

	// GetFormat returns the serialization format (yaml, json, toml)
	GetFormat() string

	// GetSectionOrder returns the order of sections as they appeared in the input file.
	GetSectionOrder() []string

	// Validate performs validation on the resume data
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

func (a *ResumeAdapter) ToResume() *Resume {
	return a.Resume
}

func (a *ResumeAdapter) GetFormat() string {
	return a.SerializationFmt
}

func (a *ResumeAdapter) GetSectionOrder() []string {
	return a.SectionOrder
}

func (a *ResumeAdapter) Validate() error {
	errors := Validate(a.Resume)
	if len(errors) > 0 {
		return fmt.Errorf("validation failed with %d errors: %v", len(errors), errors[0].Message)
	}
	return nil
}

// LoadResumeFromBytes parses resume data from raw bytes with the given format.
// Format must be one of: "yaml", "yml", "json", "toml".
func LoadResumeFromBytes(data []byte, format string) (InputData, error) {
	var resumeData Resume
	var serializationFmt string
	var sectionOrder []string

	switch strings.ToLower(format) {
	case "yaml", "yml":
		if err := UnmarshalYAMLWithContext(data, &resumeData); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
		sectionOrder = extractYAMLKeyOrder(data)
		serializationFmt = "yaml"

	case "json":
		if err := json.Unmarshal(data, &resumeData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		sectionOrder = extractJSONKeyOrder(data)
		serializationFmt = "json"

	case "toml":
		meta, err := toml.Decode(string(data), &resumeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TOML: %w", err)
		}
		sectionOrder = extractTOMLKeyOrder(meta)
		serializationFmt = "toml"

	case "md", "markdown":
		parsed, mdOrder, err := parseMarkdown(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Markdown: %w", err)
		}
		resumeData = *parsed
		sectionOrder = mdOrder
		serializationFmt = "md"

	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: yaml, yml, json, toml, md, markdown)", format)
	}

	// Fall back to default order if none was detected.
	if len(sectionOrder) == 0 {
		sectionOrder = defaultSectionOrder
	}

	// Basic validation
	if resumeData.Contact.Name == "" {
		return nil, fmt.Errorf("contact.name is required")
	}

	return &ResumeAdapter{
		Resume:           &resumeData,
		SerializationFmt: serializationFmt,
		SectionOrder:     sectionOrder,
	}, nil
}

// LoadResumeFromFile loads a resume from YAML, JSON, or TOML file.
func LoadResumeFromFile(filePath string) (InputData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	format := strings.TrimPrefix(filepath.Ext(filePath), ".")
	return LoadResumeFromBytes(data, format)
}

// extractYAMLKeyOrder parses raw YAML into a yaml.Node tree and returns the
// top-level mapping keys that correspond to renderable sections, in order.
func extractYAMLKeyOrder(data []byte) []string {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil
	}
	// root is a document node; its first child is the top-level mapping.
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil
	}
	mapping := root.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return nil
	}

	var order []string
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		key := mapping.Content[i].Value
		if sectionKeys[key] {
			order = append(order, key)
		}
	}
	return order
}

// extractJSONKeyOrder uses json.Decoder tokens to find the order of top-level
// keys that correspond to renderable sections.
func extractJSONKeyOrder(data []byte) []string {
	dec := json.NewDecoder(bytes.NewReader(data))

	// Consume opening '{'
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
				// Skip past the value for this key.
				var raw json.RawMessage
				if err := dec.Decode(&raw); err != nil {
					return order
				}
			}
		}
	}
	return order
}

// extractTOMLKeyOrder uses toml.MetaData to find the order of top-level keys
// that correspond to renderable sections.
func extractTOMLKeyOrder(meta toml.MetaData) []string {
	seen := map[string]bool{}
	var order []string
	for _, key := range meta.Keys() {
		if len(key) == 0 {
			continue
		}
		top := key[0]
		if sectionKeys[top] && !seen[top] {
			seen[top] = true
			order = append(order, top)
		}
	}
	return order
}
