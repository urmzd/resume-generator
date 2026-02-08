package resume

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// InputData represents resume data that can be validated and converted to the runtime format.
type InputData interface {
	// ToResume returns the Resume representation
	ToResume() *Resume

	// GetFormat returns the serialization format (yaml, json, toml)
	GetFormat() string

	// Validate performs validation on the resume data
	Validate() error
}

// ResumeAdapter implements InputData for Resume structures.
type ResumeAdapter struct {
	Resume           *Resume
	SerializationFmt string
}

func (a *ResumeAdapter) ToResume() *Resume {
	return a.Resume
}

func (a *ResumeAdapter) GetFormat() string {
	return a.SerializationFmt
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

	switch strings.ToLower(format) {
	case "yaml", "yml":
		if err := UnmarshalYAMLWithContext(data, &resumeData); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
		serializationFmt = "yaml"

	case "json":
		if err := json.Unmarshal(data, &resumeData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		serializationFmt = "json"

	case "toml":
		if _, err := toml.Decode(string(data), &resumeData); err != nil {
			return nil, fmt.Errorf("failed to parse TOML: %w", err)
		}
		serializationFmt = "toml"

	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: yaml, yml, json, toml)", format)
	}

	// Basic validation
	if resumeData.Contact.Name == "" {
		return nil, fmt.Errorf("contact.name is required")
	}

	return &ResumeAdapter{
		Resume:           &resumeData,
		SerializationFmt: serializationFmt,
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
