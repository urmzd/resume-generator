package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/utils"
)

func generateFilenameSlug(inputPath string) string {
	base := filepath.Base(inputPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	slug := sanitizeNameComponent(name)
	if slug == "" {
		return "resume"
	}
	return slug
}

func generateOutputBaseName(contactName string) string {
	parts := strings.Fields(contactName)
	if len(parts) == 0 {
		return "Resume"
	}
	var nameParts []string
	for _, p := range parts {
		nameParts = append(nameParts, toProperCase(p))
	}
	return strings.Join(nameParts, "_")
}

func generateRunDir(baseDir string, t time.Time) string {
	return filepath.Join(baseDir, t.Format("2006-01-02_15-04"))
}

func toProperCase(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func sanitizeNameComponent(value string) string {
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, "-", "_")

	var builder strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// sanitizeTemplateNames cleans and normalizes template names
func sanitizeTemplateNames(names []string) []string {
	var result []string
	seen := make(map[string]bool)

	for _, name := range names {
		cleaned := strings.TrimSpace(name)
		if cleaned != "" && !seen[cleaned] {
			result = append(result, cleaned)
			seen[cleaned] = true
		}
	}

	// Sort for consistent ordering
	sort.Strings(result)
	return result
}

// loadSelectedTemplates loads the specified templates or all available templates if none specified
func loadSelectedTemplates(templateNames []string) ([]*generators.Template, error) {
	if len(templateNames) == 0 {
		// Load all available templates
		allTemplates, err := generators.ListTemplates()
		if err != nil {
			return nil, fmt.Errorf("failed to list templates: %w", err)
		}

		// Convert to pointers
		var result []*generators.Template
		for i := range allTemplates {
			result = append(result, &allTemplates[i])
		}

		// Sort by name for consistent ordering
		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})

		return result, nil
	}

	// Load specified templates
	var templates []*generators.Template
	for _, name := range templateNames {
		tmpl, err := generators.LoadTemplate(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load template %s: %w", name, err)
		}
		templates = append(templates, tmpl)
	}

	return templates, nil
}

func ensureUniqueOutputPath(runDir, desiredBase, templateName, extension string) (string, error) {
	base := strings.TrimSpace(desiredBase)
	if base == "" {
		base = "Resume"
	}

	tmplSlug := strings.TrimSpace(templateName)
	if tmplSlug == "" {
		tmplSlug = "template"
	}

	ext := extension
	if ext == "" {
		ext = ".pdf"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// First attempt without numeric suffix
	candidateBase := base + "." + tmplSlug
	candidate := filepath.Join(runDir, candidateBase+ext)
	if !utils.FileExists(candidate) {
		return candidate, nil
	}

	for attempt := 2; attempt <= 9999; attempt++ {
		suffix := fmt.Sprintf("_%d", attempt)
		candidate = filepath.Join(runDir, candidateBase+suffix+ext)
		if !utils.FileExists(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("failed to find unique output filename in %s", runDir)
}
