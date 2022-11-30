package main

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/BurntSushi/toml"
	"time"
)

type Resume struct {
	Basics      Basics
	Skills      []Detail
	Experiences []Experience
	Projects    []Project
	Education   []Education
}

type Education struct {
	School   string
	Degree   string
	Suffixes []string
	Details  []Detail
	Location Location
	Dates    DateRange
}

type Basics struct {
	Name  string
	Email string
	Phone string
	Links []Link
}

type Link struct {
	Link string
}

type Location struct {
	City  string
	State string
}

type DateRange struct {
	Start time.Time
	End time.Time
}

type Experience struct {
	CompanyName string
	Title       string
	Achievement []string
	Dates       DateRange
	Location    Location
}

type Project struct {
	Name         string
	LanguageUsed string
	Details      []string
}

type Detail struct {
	Category string
	Value    string
}

type ResumeGenerator interface {
	GenerateResume(*Resume) string
}

type DefaultResume struct {
	beforeCode string
	afterCode string
}

func main() {
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	config := string(data)

	fmt.Println(config)

	if err != nil {
		panic(err)
	}
	
	var resume Resume
	_, err = toml.Decode(config, &resume)

	if err != nil {
		panic(err)
	}

	resumeJson, err := json.MarshalIndent(resume, "", "\t")
	fmt.Println(string(resumeJson))
}
