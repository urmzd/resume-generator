package main

import (
	"fmt"
	"os"

	"github.com/urmzd/generate-resumes/pkg"
	"github.com/BurntSushi/toml"
)

func main() {
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	config := string(data)

	if err != nil {
		panic(err)
	}

	var resume pkg.Resume
	_, err = toml.Decode(config, &resume)

	if err != nil {
		panic(err)
	}

	resumeBuilder := &pkg.DefaultResumeGenerator{}

	resumeBuilder.StartResume(&resume.Contact)
	resumeBuilder.AddExperiences(&resume.Experience)
	resumeBuilder.AddEducation(&resume.Education)
	resumeBuilder.AddSkills(&resume.Skills)
	if resume.Projects != nil {
		resumeBuilder.AddProjects(&resume.Projects)
	}

	resumeStr := resumeBuilder.EndResume()

	fmt.Printf("%s", resumeStr)

	//compiler := pkg.NewDefaultCompiler("xelatex")
}
