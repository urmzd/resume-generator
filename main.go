package main

import (
	"embed"
	"os"

	"github.com/urmzd/resume-generator/cmd"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var frontend embed.FS

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.EmbeddedTemplatesFS = EmbeddedTemplates
	cmd.Version = version
	cmd.Commit = commit
	cmd.BuildDate = date

	// CLI mode: any arguments → run Cobra
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	// GUI mode: no arguments → launch Wails
	app := NewApp()
	err := wails.Run(&options.App{
		Title:            "Resume Generator",
		Width:            1024,
		Height:           768,
		WindowStartState: options.Maximised,
		AssetServer: &assetserver.Options{
			Assets: frontend,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
