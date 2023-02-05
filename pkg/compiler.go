package pkg

import (
	"os"
	"log"
	"path/filepath"
	"os/exec"
)

type Compiler interface {
	LoadClasses(classes ...string)
	AddOutputFolder(string)
	Compile(string)
}

type DefaultCompiler struct {
	command string
	outputFolder string
	classes []string
}

func NewDefaultCompiler(command string) *DefaultCompiler {
	return &DefaultCompiler {
		command: command,
		outputFolder: "",
		classes: []string{},
	}

}

func (compiler *DefaultCompiler) LoadClasses(classes ...string) {
	compiler.classes = classes
}

func (compiler *DefaultCompiler) AddOutputFolder(folder string) {
	if folder == "" {
		dir, err := os.MkdirTemp("", "resume-generator");

		if err != nil {
			log.Fatal(err)
		}

		compiler.outputFolder = dir
	} else {
		compiler.outputFolder = folder
	}
}

func (compiler *DefaultCompiler) copyFile(filename string, outputFolder string) {
	data, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	newPath := filepath.Clean(filepath.Join(outputFolder, filename))

	err = os.WriteFile(newPath, data, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func (compiler *DefaultCompiler) Compile (resume string) {
	for _, class := range(compiler.classes) {
		compiler.copyFile(class, compiler.outputFolder)
	}

	outputFile , err := os.CreateTemp(compiler.outputFolder, "resume")

	if err != nil {
		log.Fatal(err)
	}

	_, err = outputFile.Write([]byte(resume));

	executable, err :=  os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	cwd := filepath.Dir(executable)

	os.Chdir(compiler.outputFolder)

	cmd := exec.Command(compiler.command, outputFile.Name())

	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	os.Chdir(cwd)
}
