package compilers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urmzd/generate-resumes/pkg/template"
	"go.uber.org/zap"
)

// XelatexCompiler is an implementation of the standard.Compiler interface
// It compiles LaTeX documents into PDFs.
type XelatexCompiler struct {
	command      string             // LaTeX compiler command (e.g., xelatex)
	outputFolder string             // Folder to store the compiled outputs
	classes      []string           // LaTeX class files to be used
	logger       *zap.SugaredLogger // Logger for logging information, warnings, and errors
}

// NewXelatexCompiler creates a new instance of DefaultCompiler with the specified command and logger.
// The command is typically a LaTeX compiler like xelatex.
func NewXelatexCompiler(command string, logger *zap.SugaredLogger) template.Compiler {
	return &XelatexCompiler{
		command:      command,
		outputFolder: "",
		classes:      []string{},
		logger:       logger,
	}
}

// LoadClasses loads LaTeX class files that will be used in the compilation.
func (compiler *XelatexCompiler) LoadClasses(classes ...string) {
	compiler.classes = classes
}

// AddOutputFolder sets the output folder for the compiled documents.
// If the folder path is not absolute, it converts it to an absolute path.
// If the folder does not exist, it creates it.
func (compiler *XelatexCompiler) AddOutputFolder(folder string) {
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
func (compiler *XelatexCompiler) Compile(resume string, resumeName string) {
	// Copy class files to the output directory
	for _, class := range compiler.classes {
		compiler.copyFile(class, compiler.outputFolder)
	}

	// Create and write the LaTeX document
	outputFilePath := filepath.Join(compiler.outputFolder, fmt.Sprintf("%s.tex", resumeName))
	err := os.WriteFile(outputFilePath, []byte(resume), 0644)
	if err != nil {
		compiler.logger.Fatal("Error creating LaTeX file:", err)
	}

	// Compile the LaTeX document
	compiler.executeLaTeXCommand(outputFilePath)
}

// copyFile copies a file from sourceFilePath to the outputFolder.
func (compiler *XelatexCompiler) copyFile(sourceFilePath, outputFolder string) {
	sourceAbsPath, err := filepath.Abs(sourceFilePath)
	if err != nil {
		compiler.logger.Fatal("Invalid source file path:", sourceFilePath)
	}

	data, err := os.ReadFile(sourceAbsPath)
	if err != nil {
		compiler.logger.Fatal("Unable to read source file:", sourceAbsPath)
	}

	destinationPath := filepath.Join(outputFolder, filepath.Base(sourceAbsPath))
	err = os.WriteFile(destinationPath, data, 0644)
	if err != nil {
		compiler.logger.Fatal("Error writing file to destination:", destinationPath)
	}
}

// executeLaTeXCommand runs the LaTeX compiler on the provided file.
func (compiler *XelatexCompiler) executeLaTeXCommand(filePath string) {
	cmd := exec.Command(compiler.command, filePath)
	cmd.Dir = compiler.outputFolder

	// Create a buffer to capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		// Log the error along with the stderr output
		compiler.logger.Fatal("LaTeX compilation error: ", err, "\nStandard Error: ", stderr.String())
	}
}
