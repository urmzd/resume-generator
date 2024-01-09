package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/urmzd/generate-resumes/pkg"
	"go.uber.org/zap"
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

		compiler := pkg.NewDefaultCompiler("xelatex", sugar)
		compiler.AddOutputFolder(OutputFolder)
		compiler.LoadClasses(ClassFiles...)

		// append timestamp to filename
		resumeFileName := "resume"
		now := time.Now()
		nowSecs := now.Unix()
		resumeFileNameFull := fmt.Sprintf("%s-%d", resumeFileName, nowSecs)

		compiler.Compile(resumeStr, resumeFileNameFull)
	},
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
