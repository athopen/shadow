package cmd

import (
	"github.com/andreaspenz/shadow/internal/io"
	"github.com/andreaspenz/shadow/internal/project"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/symfony-cli/console"
)

func newShowCommand() *console.Command {
	return &console.Command{
		Name:   "show",
		Usage:  "Show modules",
		Action: adoptActionFunc(showAction, true),
	}
}

func showAction(_ *console.Context, prj *project.Project) error {
	if len(prj.ShadowModules) == 0 {
		return errors.Errorf(`No shadow modules found at "%s"`, prj.ShadowDir)
	}

	table := tablewriter.NewWriter(io.Out())
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{
		io.Format("<header>Name</>"),
		io.Format("<header>Directory</>"),
	})

	for _, module := range prj.ShadowModules {
		table.Append([]string{
			module.Name,
			module.ModuleDir,
		})
	}

	table.Render()

	io.Write("<info>Found %d shadow module(s)</info>", len(prj.ShadowModules))

	return nil
}
