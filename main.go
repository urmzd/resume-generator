package main

import "fmt"

/**
  The pipeline for construction will go as follows:

  CreateResume(Basics) -> AddExperience -> AddSkill -> AddProject -> AddContact
*/


type Resume struct {
  experiences []Experience
  projects []Project
  skills Skills
}

type Basics struct {

} 

type Experience struct {

}

type Project struct {

}

type Skills struct {

}

func CreateResume(basics *Basics) Resume  {

}

func AddProject(*resume Resume, project Project) *Resume {

}

func AddExperience(*resume Resume, experience Experience) *Resume {

}

func AddSkills(*resume Resume, skill Skill) *Resume {

}

func main() {
  fmt.Println("Hello, World!")
}
