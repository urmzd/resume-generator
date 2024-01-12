package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/urmzd/generate-resumes/pkg/default_impl"
	"github.com/urmzd/generate-resumes/pkg/standard"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"encoding/json"
)

var rootCmd = &cobra.Command{
	Use:   "generate-resumes CONFIG",
	Short: "Generate beautiful LaTex resumes with one command.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewDevelopment()
		sugar := logger.Sugar()

		filename := args[0]

		data, err := os.ReadFile(filename)
		config := string(data)

		if err != nil {
			panic(err)
		}

		// If resume is toml, use toml.Decode else if resume is yaml, use yaml.Decode
		// else if resume is json, use json.Decode else panic.
		var resume standard.Resume
		if strings.HasSuffix(filename, ".toml") {
			_, err = toml.Decode(config, &resume)
		} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
			decoder := yaml.NewDecoder(strings.NewReader(config))
			err = decoder.Decode(&resume)
		} else if strings.HasSuffix(filename, ".json") {
			decoder := json.NewDecoder(strings.NewReader(config))
			err = decoder.Decode(&resume)
		} else {
			panic("Unknown file type.")
		}

		if err != nil {
			panic(err)
		}

		resumeBuilder := &default_impl.DefaultResumeGenerator{}

		resumeBuilder.StartResume(&resume.Contact)
		resumeBuilder.AddExperiences(&resume.Experience)
		resumeBuilder.AddEducation(&resume.Education)
		resumeBuilder.AddSkills(&resume.Skills)

		if resume.Projects != nil {
			resumeBuilder.AddProjects(&resume.Projects)
		}

		resumeStr := resumeBuilder.EndResume()

		compiler := default_impl.NewDefaultCompiler("xelatex", sugar)
		compiler.AddOutputFolder(OutputFolder)
		compiler.LoadClasses(ClassFiles...)

		contactName := strings.ReplaceAll(resume.Contact.Name, " ", "_")
		timestamp := time.Now().Format("20060102")
		versionSuffix := getVersionSuffix(contactName, OutputFolder)
		resumeFileName := fmt.Sprintf("%s_%s%s", contactName, timestamp, versionSuffix)

		compiler.Compile(resumeStr, resumeFileName)
	},
}

// getVersionSuffix checks if a file with a similar name already exists in the output folder
// and generates a version suffix if needed.
func getVersionSuffix(baseName, outputFolder string) string {
	pattern := filepath.Join(outputFolder, baseName+"_*")
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		return ""
	}

	return fmt.Sprintf("_v%d", len(files)+1)
}

var TemplateFile string
var ClassFiles []string
var OutputFolder string

func initRootCmd() {
	rootCmd.PersistentFlags().StringArrayVarP(
		&ClassFiles,
		"classes",
		"c",
		[]string{"./assets/templates/default.cls"},
		"Define the style classes that can be used.",
	)

	rootCmd.PersistentFlags().StringVarP(&OutputFolder, "outputFolder", "o", "", "Defines the location to output the compiled files.")
}

func Execute() error {
	initRootCmd()
	return rootCmd.Execute()
}
