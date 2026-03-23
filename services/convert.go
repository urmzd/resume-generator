package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urmzd/incipit/ai"
	"github.com/urmzd/incipit/utils"
)

// ConvertAndLoadResume converts a plain text file to structured JSON via LLM,
// saves the JSON alongside the original file, and loads it through the standard pipeline.
// Returns the resume data, the path to the generated JSON file, and any error.
func ConvertAndLoadResume(ctx context.Context, txtPath string, opts ai.ProviderOptions) (*ResumeData, string, error) {
	resolved, err := utils.ResolvePath(txtPath)
	if err != nil {
		return nil, "", fmt.Errorf("error resolving path: %w", err)
	}
	if !utils.FileExists(resolved) {
		return nil, "", fmt.Errorf("file does not exist: %s", resolved)
	}

	textBytes, err := os.ReadFile(resolved)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file: %w", err)
	}

	result, err := ai.Create(ctx, string(textBytes), opts)
	if err != nil {
		return nil, "", fmt.Errorf("conversion failed: %w", err)
	}

	dir := filepath.Dir(resolved)
	base := strings.TrimSuffix(filepath.Base(resolved), filepath.Ext(resolved))
	jsonPath := filepath.Join(dir, base+".json")

	if err := os.WriteFile(jsonPath, []byte(result.JSON+"\n"), 0644); err != nil {
		return nil, "", fmt.Errorf("failed to write JSON file: %w", err)
	}

	data, err := LoadResume(jsonPath)
	if err != nil {
		return nil, jsonPath, fmt.Errorf("generated JSON has validation errors (saved to %s for manual editing): %w", jsonPath, err)
	}

	return data, jsonPath, nil
}
