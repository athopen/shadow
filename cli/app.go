package cli

import (
	"fmt"
	"github.com/andreaspenz/shadow/cli/cmd"
	"github.com/symfony-cli/console"
	"time"
)

var (
	helpTemplate = `<info>
 _______ _     _ _______ ______   _____  _  _  _
 |______ |_____| |_____| |     \ |     | |  |  |
 ______| |     | |     | |_____/ |_____| |__|__|
</>

<info>{{.Name}}</>{{if .Version}} version <comment>{{.Version}}</>{{end}}{{if .Copyright}} {{.Copyright}}{{end}}

{{.Usage}}

<comment>Usage</>:
  {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} <command> [command options]{{end}} [arguments...]{{if .Description}}

{{.Description}}{{end}}{{if .VisibleFlags}}

<comment>Global options:</>
  {{range $index, $option := .VisibleFlags}}{{if $index}}
  {{end}}{{$option}}{{end}}{{end}}{{if .VisibleCommands}}

<comment>Available commands:</>{{range .VisibleCategories}}{{if .Name}}
 <comment>{{.Name}}</>{{"\t"}}{{end}}{{range .VisibleCommands}}
  <info>{{join .Names ", "}}</>{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}
`
)

func NewApplication(version string, channel string, buildDate string) *console.Application {
	return &console.Application{
		Name:      "Shadow",
		Copyright: fmt.Sprintf("(c) %d <info>Andreas Penz <andreas.penz.1989@gmail.com></>", time.Now().Year()),
		Usage:     "A tool to organize local spryker modules.",
		Commands:  cmd.GetApplicationCommands(),
		Flags:     cmd.GetApplicationFlags(),
		Action: func(ctx *console.Context) error {
			console.HelpPrinter(ctx.App.Writer, helpTemplate, ctx.App)
			return nil
		},
		Version:   version,
		Channel:   channel,
		BuildDate: buildDate,
	}
}
