package compilers

import (
	"fmt"

	"go.uber.org/zap"
)

// Engine is the unified interface for all PDF compilation backends.
// Every template type (HTML, LaTeX, Markdown) compiles content to a
// final output file (typically PDF) through an Engine.
type Engine interface {
	// Name returns a human-readable identifier for this engine (e.g. "rod", "xelatex").
	Name() string
	// Compile transforms rendered content into a final output file at outputPath.
	Compile(content string, outputPath string) error
}

// EngineConfig holds the per-type engine selections read from the config file.
type EngineConfig struct {
	HTML     string // e.g. "rod"
	LaTeX    string // e.g. "xelatex", "pdflatex", "lualatex", "auto"
	Markdown string // e.g. "rod"
}

// NewEngine creates an Engine for the given template type using the provided config.
// templateDir is only used for LaTeX engines (to copy class files).
func NewEngine(templateType string, cfg EngineConfig, templateDir string, logger *zap.SugaredLogger) (Engine, error) {
	switch templateType {
	case "html":
		return NewHTMLEngine(cfg.HTML, logger)
	case "latex":
		return NewLaTeXEngine(cfg.LaTeX, templateDir, logger)
	case "markdown":
		return NewMarkdownEngine(cfg.Markdown, logger)
	default:
		return nil, fmt.Errorf("no engine available for template type: %s", templateType)
	}
}

// NewHTMLEngine creates an HTML-to-PDF engine. Currently only "rod" is supported.
func NewHTMLEngine(name string, logger *zap.SugaredLogger) (Engine, error) {
	if name == "" {
		name = "rod"
	}
	switch name {
	case "rod":
		return NewRodHTMLToPDFCompiler(logger), nil
	default:
		return nil, fmt.Errorf("unknown HTML engine: %q (available: rod)", name)
	}
}

// NewLaTeXEngine creates a LaTeX-to-PDF engine.
// name can be "auto", "xelatex", "pdflatex", "lualatex", or "latex".
func NewLaTeXEngine(name string, templateDir string, logger *zap.SugaredLogger) (Engine, error) {
	if name == "" || name == "auto" {
		engine := DetectLaTeXEngine()
		if engine == "" {
			return nil, fmt.Errorf("no LaTeX engine found\n\nPlease install one of the following:\n  - TeX Live:   https://www.tug.org/texlive/\n  - MiKTeX:     https://miktex.org/\n  - MacTeX:     https://www.tug.org/mactex/ (macOS)")
		}
		name = engine
		logger.Infof("Auto-detected LaTeX engine: %s", name)
	}
	return &LaTeXEngine{
		command:     name,
		templateDir: templateDir,
		logger:      logger,
	}, nil
}

// NewMarkdownEngine creates a Markdown-to-PDF engine.
// It converts Markdown to HTML, then delegates to an HTML engine.
func NewMarkdownEngine(name string, logger *zap.SugaredLogger) (Engine, error) {
	htmlEngine, err := NewHTMLEngine(name, logger)
	if err != nil {
		return nil, fmt.Errorf("markdown engine requires an HTML engine: %w", err)
	}
	return &markdownEngine{htmlEngine: htmlEngine}, nil
}
