package services

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/urmzd/incipit/resume"
)

// GenerateSchema generates the JSON schema for the resume format and returns
// the marshalled JSON bytes.
func GenerateSchema() ([]byte, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema := reflector.Reflect(&resume.Resume{})

	schema.Title = "Resume Format (v2.0)"
	schema.Description = `Unified resume format used by the CLI.

Supported serialization formats:
- YAML (.yml, .yaml)
- JSON (.json)
- TOML (.toml)

Key features:
- Contact details with optional links and location
- Experience, education, projects, and skills sections
- Date range validation for time-based entries
- Location information with city/state/country`

	addSchemaExample(schema)

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	return schemaJSON, nil
}

func addSchemaExample(schema *jsonschema.Schema) {
	schema.Examples = []interface{}{
		map[string]interface{}{
			"contact": map[string]interface{}{
				"name":  "Jane Smith",
				"email": "jane.smith@example.com",
				"phone": "+1-555-987-6543",
				"links": []map[string]interface{}{
					{
						"uri": "https://github.com/janesmith",
					},
				},
				"location": map[string]interface{}{
					"city":    "San Francisco",
					"state":   "CA",
					"country": "USA",
				},
			},
			"experience": map[string]interface{}{
				"positions": []map[string]interface{}{
					{
						"company": "Tech Innovations Inc",
						"title":   "Senior Software Engineer",
						"highlights": []string{
							"Led team of 5 engineers in building microservices architecture",
							"Improved system performance by 60% through optimization",
						},
						"dates": map[string]interface{}{
							"start": "2021-06-01T00:00:00Z",
							"end":   "2024-01-01T00:00:00Z",
						},
						"location": map[string]interface{}{
							"city":  "San Francisco",
							"state": "CA",
						},
					},
				},
			},
			"skills": map[string]interface{}{
				"title": "Technical Skills",
				"categories": []map[string]interface{}{
					{
						"category": "Programming Languages",
						"items": []string{
							"Go",
							"Python",
						},
					},
				},
			},
			"education": map[string]interface{}{
				"title": "Education",
				"institutions": []map[string]interface{}{
					{
						"institution": "University of California, Berkeley",
						"degree": map[string]interface{}{
							"name": "Bachelor of Science in Computer Science",
						},
						"dates": map[string]interface{}{
							"start": "2013-08-01T00:00:00Z",
							"end":   "2017-05-15T00:00:00Z",
						},
						"gpa": map[string]interface{}{
							"gpa":     "3.8",
							"max_gpa": "4.0",
						},
					},
				},
			},
		},
	}
}
