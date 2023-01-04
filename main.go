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

	//fmt.Printf("%+v\n", &resume)

	resumeBuilder := &pkg.DefaultResumeGenerator{}

	resumeBuilder.StartResume(&resume.Contact)
	resumeBuilder.AddSkills(&resume.Skills)
	resumeBuilder.AddExperiences(&resume.Experience)
	resumeBuilder.AddProjects(&resume.Projects)
	resumeBuilder.AddEducation(&resume.Education)

	fmt.Println(resumeBuilder.EndResume())
}
