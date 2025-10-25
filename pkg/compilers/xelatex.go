package compilers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"

	"io"
)

// DetectLaTeXEngine tries to find an available LaTeX compiler in the system.
// It checks for engines in the following order: xelatex, pdflatex, lualatex, latex
// Returns the name of the first available engine, or an empty string if none are found.
func DetectLaTeXEngine() string {
	engines := []string{"xelatex", "pdflatex", "lualatex", "latex"}

	for _, engine := range engines {
		if _, err := exec.LookPath(engine); err == nil {
			return engine
		}
	}

	return ""
}

// GetAvailableLaTeXEngines returns a list of all available LaTeX engines on the system.
func GetAvailableLaTeXEngines() []string {
	engines := []string{"xelatex", "pdflatex", "lualatex", "latex"}
	available := []string{}

	for _, engine := range engines {
		if _, err := exec.LookPath(engine); err == nil {
			available = append(available, engine)
		}
	}

	return available
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// copyDir copies the contents of the src directory to dst.
// Does not copy subdirectories.
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		fileInfo, err := entry.Info()
		if err != nil {
			return err
		}

		if fileInfo.Mode().IsRegular() {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// LaTeXCompiler is an implementation of the standard.Compiler interface
// It compiles LaTeX documents into PDFs using various LaTeX engines.
type LaTeXCompiler struct {
	command      string             // LaTeX compiler command (e.g., xelatex, pdflatex, lualatex)
	outputFolder string             // Folder to store the compiled outputs
	classes      string             // LaTeX class files to be used
	logger       *zap.SugaredLogger // Logger for logging information, warnings, and errors
}

// NewLaTeXCompiler creates a new instance of LaTeXCompiler with the specified command and logger.
// The command can be any LaTeX compiler like xelatex, pdflatex, lualatex, etc.
func NewLaTeXCompiler(command string, logger *zap.SugaredLogger) definition.Compiler {
	return &LaTeXCompiler{
		command:      command,
		outputFolder: "",
		classes:      "",
		logger:       logger,
	}
}

// NewAutoLaTeXCompiler creates a LaTeX compiler by automatically detecting the available engine.
// Returns nil if no LaTeX engine is found.
func NewAutoLaTeXCompiler(logger *zap.SugaredLogger) (definition.Compiler, error) {
	engine := DetectLaTeXEngine()
	if engine == "" {
		return nil, fmt.Errorf("no LaTeX engine found\n\nPlease install one of the following:\n  - TeX Live:   https://www.tug.org/texlive/\n  - MiKTeX:     https://miktex.org/\n  - MacTeX:     https://www.tug.org/mactex/ (macOS)\n\nOr use Docker which includes LaTeX:\n  docker run --rm -v $(pwd):/work texlive/texlive")
	}

	logger.Infof("Auto-detected LaTeX engine: %s", engine)
	return NewLaTeXCompiler(engine, logger), nil
}

// XelatexCompiler is deprecated. Use LaTeXCompiler instead.
type XelatexCompiler = LaTeXCompiler

// NewXelatexCompiler is deprecated. Use NewLaTeXCompiler or NewAutoLaTeXCompiler instead.
func NewXelatexCompiler(command string, logger *zap.SugaredLogger) definition.Compiler {
	return NewLaTeXCompiler(command, logger)
}

// LoadClasses loads LaTeX class files that will be used in the compilation.
func (compiler *LaTeXCompiler) LoadClasses(classes string) {
	compiler.classes = classes
}

// AddOutputFolder sets the output folder for the compiled documents.
// If the folder path is not absolute, it converts it to an absolute path.
// If the folder does not exist, it creates it.
func (compiler *LaTeXCompiler) AddOutputFolder(folder string) {
	var err error
	if folder == "" {
		compiler.outputFolder, err = os.MkdirTemp("", "resume-generator")
	} else {
		compiler.outputFolder, err = filepath.Abs(folder)
		if err != nil {
			err = os.MkdirAll(compiler.outputFolder, 0755)
		}
	}

	if err != nil {
		compiler.logger.Fatal("Error setting output folder:", err)
	}
}

// Compile compiles the LaTeX document into a PDF.
// It copies necessary class files to the output directory, creates the .tex file,
// and then runs the LaTeX compiler.
func (compiler *LaTeXCompiler) Compile(resume string, resumeName string) string {
	// Copy the class files to the output folder
	copyDir(compiler.classes, compiler.outputFolder)

	// Create and write the LaTeX document
	outputFilePath := filepath.Join(compiler.outputFolder, fmt.Sprintf("%s.tex", resumeName))
	err := os.WriteFile(outputFilePath, []byte(resume), 0644)
	if err != nil {
		compiler.logger.Fatal("Error creating LaTeX file:", err)
	}

	// Compile the LaTeX document
	compiler.executeLaTeXCommand(outputFilePath)

	return outputFilePath
}

// executeLaTeXCommand runs the LaTeX compiler on the provided file.
func (compiler *LaTeXCompiler) executeLaTeXCommand(filePath string) {
	cmd := exec.Command(compiler.command, filePath)
	cmd.Dir = compiler.outputFolder

	// Create a buffer to capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		// Log the error along with the stderr output
		compiler.logger.Fatalf("LaTeX compilation error with %s: %v\nStandard Error: %s", compiler.command, err, stderr.String())
	}

	compiler.logger.Infof("Successfully compiled with %s", compiler.command)
}
