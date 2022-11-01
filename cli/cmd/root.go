package cmd

import (
	"github.com/andreaspenz/shadow/internal/filesystem"
	"github.com/andreaspenz/shadow/internal/project"
	"github.com/symfony-cli/console"
	"path"
)

var (
	projectDirFlag = &console.StringFlag{Name: "project-dir", Usage: "Specify the project directory", DefaultValue: "."}
)

func GetApplicationFlags() []console.Flag {
	return []console.Flag{
		projectDirFlag,
	}
}

func GetApplicationCommands() []*console.Command {
	return []*console.Command{
		newCleanCommand(),
		newDeployCommand().cc,
		newInitCommand().cc,
		newReverseCommand().cc,
		newShowCommand(),
	}
}

type adoptedActionFunc func(ctx *console.Context, prj *project.Project) error

func adoptActionFunc(fn adoptedActionFunc, fullLoad bool) console.ActionFunc {
	return func(ctx *console.Context) (err error) {
		prj, err := project.LoadProject(project.Descriptor{
			Fs:         filesystem.NewOsFs(),
			ProjectDir: path.Clean(ctx.String(projectDirFlag.Name)),
		}, fullLoad)

		if err != nil {
			return err
		}

		return fn(ctx, prj)
	}
}
