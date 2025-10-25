package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/definition"
)

var (
	SchemaFormat string
	SchemaOutput string
)

func initSchemaCmd() {
	schemaCmd.AddCommand(schemaGenerateCmd)
	schemaCmd.AddCommand(schemaListCmd)
	rootCmd.AddCommand(schemaCmd)

	schemaGenerateCmd.Flags().StringVarP(&SchemaFormat, "format", "f", "all", "Schema format to generate (legacy, enhanced, json-resume, all)")
	schemaGenerateCmd.Flags().StringVarP(&SchemaOutput, "output", "o", "", "Output directory for schema files (defaults to stdout)")
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate JSON schemas for resume formats",
	Long: `Generate JSON schemas that can be used by LLMs, validators, and tools
to understand and work with resume formats.`,
}

var schemaGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate JSON schema files",
	Long: `Generate JSON Schema files for resume formats. These schemas can be used by:
- LLMs to understand the structure and generate/enhance resumes
- Validation tools to verify resume data
- IDEs for autocompletion and validation
- Documentation generation`,
	Run: func(cmd *cobra.Command, args []string) {
		formats := []string{}
		switch SchemaFormat {
		case "all":
			formats = []string{"legacy", "enhanced", "json-resume"}
		case "legacy":
			formats = []string{"legacy"}
		case "enhanced":
			formats = []string{"enhanced"}
		case "json-resume":
			formats = []string{"json-resume"}
		default:
			fmt.Printf("Unknown format: %s\n", SchemaFormat)
			fmt.Println("Valid formats: legacy, enhanced, json-resume, all")
			return
		}

		for _, format := range formats {
			if err := generateSchema(format); err != nil {
				fmt.Printf("Error generating %s schema: %v\n", format, err)
			}
		}
	},
}

var schemaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available schema formats",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available Resume Formats:")
		fmt.Println()

		fmt.Println("📋 Legacy Format (YAML/JSON/TOML)")
		fmt.Println("  - Traditional resume format")
		fmt.Println("  - Simple structure with basic fields")
		fmt.Println("  - Backward compatible")
		fmt.Println("  - Generate schema: resume-generator schema generate -f legacy")
		fmt.Println()

		fmt.Println("✨ Enhanced Format (v2.0)")
		fmt.Println("  - Modern format with advanced features")
		fmt.Println("  - Supports ordering, visibility controls, metrics")
		fmt.Println("  - Template embedding and theming")
		fmt.Println("  - Generate schema: resume-generator schema generate -f enhanced")
		fmt.Println()

		fmt.Println("🌐 JSON Resume Format")
		fmt.Println("  - Community standard format (jsonresume.org)")
		fmt.Println("  - Wide tool compatibility")
		fmt.Println("  - Standard schema format")
		fmt.Println("  - Generate schema: resume-generator schema generate -f json-resume")
		fmt.Println()

		fmt.Println("Usage:")
		fmt.Println("  # Generate all schemas")
		fmt.Println("  resume-generator schema generate")
		fmt.Println()
		fmt.Println("  # Generate specific format")
		fmt.Println("  resume-generator schema generate -f enhanced")
		fmt.Println()
		fmt.Println("  # Save to directory")
		fmt.Println("  resume-generator schema generate -o ./schemas")
	},
}

func generateSchema(format string) error {
	var schema *jsonschema.Schema
	var schemaType interface{}
	var filename string

	switch format {
	case "legacy":
		schemaType = &definition.Resume{}
		filename = "resume-legacy.schema.json"
	case "enhanced":
		schemaType = &definition.EnhancedResume{}
		filename = "resume-enhanced.schema.json"
	case "json-resume":
		schemaType = &definition.JSONResume{}
		filename = "resume-jsonresume.schema.json"
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	// Generate schema
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema = reflector.Reflect(schemaType)

	// Add metadata
	schema.Title = fmt.Sprintf("Resume %s Format", format)
	schema.Description = getSchemaDescription(format)

	// Add examples
	addSchemaExamples(schema, format)

	// Marshal to JSON
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Output
	if SchemaOutput != "" {
		// Save to file
		if err := os.MkdirAll(SchemaOutput, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		outputPath := fmt.Sprintf("%s/%s", SchemaOutput, filename)
		if err := os.WriteFile(outputPath, schemaJSON, 0644); err != nil {
			return fmt.Errorf("failed to write schema file: %w", err)
		}

		fmt.Printf("✓ Generated %s schema: %s\n", format, outputPath)
	} else {
		// Print to stdout
		fmt.Printf("\n=== %s Schema ===\n", format)
		fmt.Println(string(schemaJSON))
	}

	return nil
}

func getSchemaDescription(format string) string {
	descriptions := map[string]string{
		"legacy": `Traditional resume format with simple structure. Supports contact information,
experience, education, skills, and projects. Compatible with YAML, JSON, and TOML.

Use this format for:
- Simple, straightforward resumes
- Backward compatibility
- Quick resume generation

Fields support:
- Contact: name, email, phone, links
- Experience: company, title, achievements/description, dates, location
- Education: school, degree, suffixes, details, dates
- Skills: category-value pairs
- Projects: name, language, description, link`,

		"enhanced": `Modern resume format (v2.0) with advanced features including ordering,
visibility controls, metrics, and template embedding.

Use this format for:
- Complex, detailed resumes
- Custom ordering of sections
- Visibility control over individual fields
- Quantifiable achievements with metrics
- Multi-language support
- Theme customization

Advanced features:
- Section ordering and visibility
- Achievement metrics and quantification
- Template embedding
- Multiple output formats
- Custom themes
- Location and date range validation`,

		"json-resume": `Standard JSON Resume format following the jsonresume.org specification.
Widely supported by resume tools and websites.

Use this format for:
- Maximum compatibility with other tools
- Publishing to JSON Resume platforms
- Standard resume structure
- Web-based resume hosting

Follows the official JSON Resume schema with support for:
- Basics (contact, summary, profiles)
- Work experience
- Volunteer experience
- Education
- Awards and publications
- Skills and languages
- Interests and references`,
	}

	return descriptions[format]
}

func addSchemaExamples(schema *jsonschema.Schema, format string) {
	// Add format-specific examples
	switch format {
	case "legacy":
		schema.Examples = []interface{}{
			map[string]interface{}{
				"contact": map[string]interface{}{
					"name":  "John Doe",
					"email": "john.doe@example.com",
					"phone": "+1-555-123-4567",
					"links": []map[string]interface{}{
						{"link": "github.com/johndoe"},
						{"link": "linkedin.com/in/johndoe"},
					},
				},
				"experience": []map[string]interface{}{
					{
						"company": "Tech Corp",
						"title":   "Software Engineer",
						"achievements": []string{
							"Developed microservices handling 1M+ requests/day",
							"Reduced API latency by 40% through optimization",
						},
						"dates": map[string]interface{}{
							"start": "2020-01-01",
							"end":   "2023-12-31",
						},
					},
				},
				"skills": []map[string]interface{}{
					{
						"category": "Languages",
						"value":    "Python, Go, JavaScript, TypeScript",
					},
				},
			},
		}

	case "enhanced":
		schema.Examples = []interface{}{
			map[string]interface{}{
				"meta": map[string]interface{}{
					"version": "2.0",
					"output": map[string]interface{}{
						"formats": []string{"pdf", "html"},
					},
					"theme": "modern",
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
					"visibility": map[string]interface{}{
						"showPhone":    true,
						"showEmail":    true,
						"showLocation": true,
					},
				},
				"experience": map[string]interface{}{
					"order": 2,
					"title": "Professional Experience",
					"positions": []map[string]interface{}{
						{
							"order":   1,
							"company": "Innovative Solutions Inc",
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
						},
					},
				},
			},
		}

	case "json-resume":
		schema.Examples = []interface{}{
			map[string]interface{}{
				"basics": map[string]interface{}{
					"name":    "Richard Hendricks",
					"label":   "Programmer",
					"email":   "richard.hendricks@example.com",
					"phone":   "+1-555-123-4567",
					"url":     "https://richardhendricks.example.com",
					"summary": "Richard is a talented developer with a strong background in compression algorithms.",
					"location": map[string]interface{}{
						"city":        "San Francisco",
						"countryCode": "US",
						"region":      "California",
					},
					"profiles": []map[string]interface{}{
						{
							"network":  "Twitter",
							"username": "richard",
							"url":      "https://twitter.com/richard",
						},
					},
				},
				"work": []map[string]interface{}{
					{
						"name":      "Pied Piper",
						"position":  "CEO/President",
						"startDate": "2013-12-01",
						"endDate":   "2014-12-01",
						"summary":   "Pied Piper is a multi-platform technology based on a proprietary universal compression algorithm.",
						"highlights": []string{
							"Built an algorithm for artist to detect if their music was violating copy right infringement laws",
							"Developed a web application that determines how efficient code is",
						},
					},
				},
			},
		}
	}
}
