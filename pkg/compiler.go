package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"go.uber.org/zap"
)

type Compiler interface {
	LoadClasses(classes ...string)
	AddOutputFolder(string)
	Compile(string, string)
}

type DefaultCompiler struct {
	command string
	outputFolder string
	classes []string
	logger *zap.SugaredLogger
}

func NewDefaultCompiler(command string, logger *zap.SugaredLogger) Compiler {
	return &DefaultCompiler {
		command: command,
		outputFolder: "",
		classes: []string{},
		logger: logger,
	}

}

func (compiler *DefaultCompiler) LoadClasses(classes ...string) {
	compiler.classes = classes
}

func (compiler *DefaultCompiler) AddOutputFolder(folder string) {
	if folder == "" {
		dir, err := os.MkdirTemp("", "resume-generator");

		if err != nil {
			compiler.logger.Fatal(err)
		}

		compiler.outputFolder = dir
	} else {
		compiler.outputFolder = folder
	}
}

func (compiler *DefaultCompiler) copyFile(sourceFilePath string, outputFolder string) {
	absFilepath, err := filepath.Abs(sourceFilePath)

	if err != nil {
		compiler.logger.Fatal(err)
	}

	data, err := os.ReadFile(absFilepath)

	if err != nil {
		compiler.logger.Error(err)
	}

	fileName := filepath.Base(absFilepath)
	newPath := filepath.Clean(filepath.Join(outputFolder, fileName))

	compiler.logger.Info(absFilepath, newPath)

	err = os.WriteFile(newPath, data, 0644)

	if err != nil {
		compiler.logger.Error(err)
	}
}

func (compiler *DefaultCompiler) Compile (resume string, resume_name string) {
	// Copy style classes over to temporary directory.
	compiler.logger.Info(compiler.classes)

	for _, class := range(compiler.classes) {
		compiler.copyFile(class, compiler.outputFolder)
	}

	// Create the resume tex file.
	outputFileName := fmt.Sprintf("%s.tex", resume_name);
	outputFile , err := os.Create(filepath.Join(compiler.outputFolder, outputFileName));

	if err != nil {
		compiler.logger.Fatal(err)
	}

	// Copy the code over.
	_, err = outputFile.Write([]byte(resume));

	executable, err := os.Executable() 

	if err != nil {
		compiler.logger.Fatal(err)
	}

	cwd := filepath.Dir(executable)

	// Switch to temporary directory.
	os.Chdir(compiler.outputFolder)

	// Create the compilation command.
	compiler.logger.Info(outputFile.Name())
	cmd := exec.Command(compiler.command, outputFile.Name())

	// Run the command.
	err = cmd.Run()

	if err != nil {
		compiler.logger.Fatal(err)
	}

	// Clean
	dirFiles, err := os.ReadDir("./")

	if err != nil {
		compiler.logger.Fatal(err)
	}

	for _, file := range dirFiles {
		baseName := file.Name()
		baseExt := filepath.Ext(baseName)

		compiler.logger.Info(baseName, baseExt)

		if baseExt != ".pdf" {
			os.Remove(filepath.Clean(filepath.Join("./", baseName)))
		}
	}

	compiler.logger.Infow("Compilation completed.", "outputFolder", compiler.outputFolder)

	defer os.Chdir(cwd)
}
