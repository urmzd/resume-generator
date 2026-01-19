package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/resume"
)

var SchemaOutput string

func initSchemaCmd() {
	rootCmd.AddCommand(schemaCmd)
	schemaCmd.Flags().StringVarP(&SchemaOutput, "output", "o", "", "Output file path (defaults to stdout for easy piping)")
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Output JSON schema for resume input format",
	Long: `Output JSON Schema for resume input validation. Outputs to stdout by default,
making it easy to pipe to clipboard or save to a file.

The resume generator uses a single resume format (v2.0) that supports:
- Contact details with links and locations
- Experience, education, projects, and skills sections
- Date ranges for time-based entries
- Multiple serialization formats (YAML, JSON, TOML)

Examples:
  # Output schema to stdout
  resume-generator schema

  # Copy schema to clipboard (macOS)
  resume-generator schema | pbcopy

  # Copy schema to clipboard (Linux with xclip)
  resume-generator schema | xclip -selection clipboard

  # Save schema to file
  resume-generator schema -o resume-schema.json

  # Use with validation tools
  resume-generator schema | jq '.'`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateSchema(); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating schema: %v\n", err)
			os.Exit(1)
		}
	},
}

func generateSchema() error {
	// Generate schema for the resume format
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema := reflector.Reflect(&resume.Resume{})

	// Add metadata
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

	// Add example
	addSchemaExample(schema)

	// Marshal to JSON
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Output
	if SchemaOutput != "" {
		// Save to file
		if err := os.WriteFile(SchemaOutput, schemaJSON, 0644); err != nil {
			return fmt.Errorf("failed to write schema file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Schema written to: %s\n", SchemaOutput)
	} else {
		// Print to stdout
		fmt.Println(string(schemaJSON))
	}

	return nil
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
