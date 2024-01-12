package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/urmzd/generate-resumes/pkg/compilers"
	"github.com/urmzd/generate-resumes/pkg/generators"
	"github.com/urmzd/generate-resumes/pkg/template"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"encoding/json"
)

var (
	TemplateFile  string
	ClassFiles    []string
	OutputFolder  string
	TemplateTypes []string
	KeepTex       bool
)

func initRootCmd() {
	rootCmd.PersistentFlags().StringArrayVarP(&ClassFiles, "classes", "c", []string{"./assets/templates/default.cls"}, "Define the style classes that can be used.")
	rootCmd.PersistentFlags().StringVarP(&OutputFolder, "output-folder", "o", "", "Defines the location to output the compiled files.")
	rootCmd.PersistentFlags().StringSliceVarP(&TemplateTypes, "templates", "t", []string{"professional", "creative", "technical", "base"}, "Specify which resume templates to use.")
	rootCmd.PersistentFlags().BoolVarP(&KeepTex, "keep-tex", "k", false, "Keep .tex files after compilation")
}

var rootCmd = &cobra.Command{
	Use:   "generate-resumes CONFIG",
	Short: "Generate beautiful LaTeX resumes with one command.",
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

		var resume template.Resume
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

		for _, templateType := range TemplateTypes {
			var resumeBuilder template.ResumeGenerator

			switch templateType {
			case "base":
				resumeBuilder = generators.NewBaseResumeGenerator()
			case "professional":
				continue
				// resumeBuilder = generators.NewProfessionalResumeGenerator()
			case "creative":
				// Skip for now
				continue
				// resumeBuilder = generators.NewCreativeResumeGenerator()
			case "technical":
				continue
				// resumeBuilder = generators.NewTechnicalResumeGenerator()
			default:
				sugar.Fatalf("Unknown template type: %s", templateType)
			}

			resumeBuilder.StartResume(&resume.Contact)
			resumeBuilder.AddExperiences(&resume.Experience)
			resumeBuilder.AddEducation(&resume.Education)
			resumeBuilder.AddSkills(&resume.Skills)

			if resume.Projects != nil {
				resumeBuilder.AddProjects(&resume.Projects)
			}

			resumeStr := resumeBuilder.EndResume()

			compiler := compilers.NewXelatexCompiler("xelatex", sugar)
			compiler.AddOutputFolder(OutputFolder)
			compiler.LoadClasses(ClassFiles...)

			contactName := strings.ReplaceAll(resume.Contact.Name, " ", "_")
			timestamp := time.Now().Format("20060102")
			versionSuffix := getVersionSuffix(contactName, OutputFolder)
			resumeFileName := fmt.Sprintf("%s_%s_%s%s", contactName, timestamp, templateType, versionSuffix)

			compiler.Compile(resumeStr, resumeFileName)

			if KeepTex {
				cleanArtifacts(sugar, ".tex")
			} else {
				cleanArtifacts(sugar)
			}
		}
	},
}

func getVersionSuffix(baseName, outputFolder string) string {
	pattern := filepath.Join(outputFolder, baseName+"_*")
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		return ""
	}

	return fmt.Sprintf("_v%d", len(files)+1)
}

func cleanArtifacts(logger *zap.SugaredLogger, keepExtensions ...string) {
	files, err := os.ReadDir(OutputFolder)
	if err != nil {
		logger.Fatal("Error reading output folder:", err)
	}

	// Create a map for quick lookup of extensions to keep
	keep := make(map[string]bool)
	for _, ext := range keepExtensions {
		keep[ext] = true
	}
	// Always keep PDF files
	keep[".pdf"] = true

	for _, file := range files {
		if _, ok := keep[filepath.Ext(file.Name())]; !ok {
			// If the file's extension is not in the keep list, remove it
			err := os.Remove(filepath.Join(OutputFolder, file.Name()))
			if err != nil {
				logger.Error("Error removing file:", err)
			}
		}
	}
}

func Execute() error {
	initRootCmd()
	return rootCmd.Execute()
}
