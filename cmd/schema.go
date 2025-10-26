package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/definition"
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
- Section ordering and visibility controls
- Quantifiable achievements with metrics
- Multiple serialization formats (YAML, JSON, TOML)
- Advanced features like certifications, publications, and languages

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
	schema := reflector.Reflect(&definition.Resume{})

	// Add metadata
	schema.Title = "Resume Format (v2.0)"
	schema.Description = `Modern resume format with advanced features including ordering,
visibility controls, metrics, and comprehensive professional information.

Supported serialization formats:
- YAML (.yml, .yaml)
- JSON (.json)
- TOML (.toml)

Key features:
- Section ordering and visibility controls
- Quantifiable achievements with metrics
- Certifications, publications, and languages support
- Date range validation
- Location information with multiple levels of detail
- Template embedding and theming`

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
			"meta": map[string]interface{}{
				"version": "2.0",
			},
			"contact": map[string]interface{}{
				"order": 1,
				"name":  "Jane Smith",
				"title": "Senior Software Engineer",
				"email": "jane.smith@example.com",
				"phone": "+1-555-987-6543",
				"links": []map[string]interface{}{
					{
						"order": 1,
						"text":  "GitHub",
						"url":   "https://github.com/janesmith",
						"type":  "github",
					},
				},
				"location": map[string]interface{}{
					"city":    "San Francisco",
					"state":   "CA",
					"country": "USA",
				},
				"summary": "Experienced software engineer specializing in distributed systems and cloud architecture.",
			},
			"experience": map[string]interface{}{
				"order": 2,
				"title": "Professional Experience",
				"positions": []map[string]interface{}{
					{
						"order":   1,
						"company": "Tech Innovations Inc",
						"title":   "Senior Software Engineer",
						"highlights": []string{
							"Led team of 5 engineers in building microservices architecture",
							"Improved system performance by 60% through optimization",
						},
						"achievements": []map[string]interface{}{
							{
								"order":       1,
								"description": "Reduced deployment time by 75%",
								"metric": map[string]interface{}{
									"name":   "deployment_time",
									"before": 40,
									"after":  10,
									"unit":   "minutes",
								},
							},
						},
						"dates": map[string]interface{}{
							"start":   "2021-06-01T00:00:00Z",
							"current": true,
						},
						"location": map[string]interface{}{
							"city":  "San Francisco",
							"state": "CA",
						},
					},
				},
			},
			"skills": map[string]interface{}{
				"order": 3,
				"title": "Technical Skills",
				"categories": []map[string]interface{}{
					{
						"order": 1,
						"name":  "Programming Languages",
						"items": []map[string]interface{}{
							{
								"name":              "Go",
								"level":             "Expert",
								"yearsOfExperience": 5,
							},
							{
								"name":              "Python",
								"level":             "Advanced",
								"yearsOfExperience": 7,
							},
						},
					},
				},
			},
			"education": map[string]interface{}{
				"order": 4,
				"title": "Education",
				"institutions": []map[string]interface{}{
					{
						"order":       1,
						"institution": "University of California, Berkeley",
						"degree":      "Bachelor of Science in Computer Science",
						"dates": map[string]interface{}{
							"start": "2013-08-01T00:00:00Z",
							"end":   "2017-05-15T00:00:00Z",
						},
						"gpa": "3.8",
					},
				},
			},
		},
	}
}
