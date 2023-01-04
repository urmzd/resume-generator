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

func (compiler *XeLaTeXCompiler) Compile (folder string) {

}
