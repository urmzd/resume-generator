package pkg

import (
	"os"
	"log"
)

type Compiler interface {
	LoadClasses(classes ...string)
	AddOutputFolder(string)
	Compile(string)
}

type XeLaTeXCompiler struct {
	command string
	outputFolder string
	classes []string
}

func (compiler *XeLaTeXCompiler) LoadClasses(classes ...string) {
	compiler.classes = classes
}

func (compiler *XeLaTeXCompiler) AddOutputFolder(folder string) {
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

func (compiler *XeLaTeXCompiler) copyFile(filename string, outputFolder string) {
	data, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Figure out which permissions to use, default 644)
	err = os.WriteFile(filename, data, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func (compiler *XeLaTeXCompiler) Compile (resume string) {

	for _, class := range(compiler.classes) {
		compiler.copyFile(class, compiler.outputFolder)
	}

	outputFile , err := os.CreateTemp(compiler.outputFolder, "resume")

	if err != nil {
		log.Fatal(err)
	}

	_, err = outputFile.Write([]byte(resume));

	// TODO: Get current working directory.

	//cmd := `pdflatex <resume_file>`

	// TODO: Run cmd

	// 

}
