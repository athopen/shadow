package cmd

import (
	"github.com/andreaspenz/shadow/internal/filesystem"
	"github.com/andreaspenz/shadow/internal/io"
	"github.com/andreaspenz/shadow/internal/project"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/symfony-cli/console"
	"path/filepath"
)

var (
	forceFlag = &console.BoolFlag{Name: "force", Usage: "Overwrite existing files, directories and symlinks"}
	copyFlag  = &console.BoolFlag{Name: "copy", Usage: "Copy and paste instead of symlink"}
)

type deployCommand struct {
	cc *console.Command
}

func newDeployCommand() *deployCommand {
	cmd := &deployCommand{}

	cmd.cc = &console.Command{
		Name:   "deploy",
		Usage:  "Deploy modules",
		Flags:  []console.Flag{forceFlag, copyFlag},
		Action: adoptActionFunc(cmd.deployAction, true),
	}

	return cmd
}

func (d *deployCommand) deployAction(ctx *console.Context, prj *project.Project) error {
	if len(prj.ShadowModules) == 0 {
		return errors.Errorf(`No modules found at "%s"`, prj.ShadowDir)
	}

	forceMode := ctx.Bool(forceFlag.Name)
	copyMode := ctx.Bool(copyFlag.Name)

	var deployed int
	for _, module := range prj.ShadowModules {
		proceeded, err := d.deployModule(prj, module, forceMode, copyMode)

		if err != nil {
			return err
		}

		if proceeded {
			deployed++
		}
	}

	if deployed == 0 {
		io.Write("<info>No new links found to deploy</info>")
		return nil
	}

	io.Write("<info>Deployment of %d module(s) succeeded</info>", deployed)

	return nil
}

func (d *deployCommand) deployModule(prj *project.Project, module *project.ShadowModule, forceMode bool, copyMode bool) (bool, error) {
	var deployed bool
	for from, to := range module.Links {
		applied, err := d.deployLink(prj, forceMode, copyMode, from, to)

		if err != nil {
			return false, err
		}

		if applied {
			deployed = true
		}
	}

	return deployed, nil
}

func (d *deployCommand) deployLink(prj *project.Project, forceMode bool, copyMode bool, from string, to string) (bool, error) {
	// skip if to already exists and not in forceMode
	if exists, _ := filesystem.Exists(prj.Fs, to); exists && !forceMode {
		return false, nil
	}

	// remove to if it already exists
	if err := prj.Fs.RemoveAll(to); err != nil {
		return false, err
	}

	// create parent directories
	if err := prj.Fs.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return false, err
	}

	if copyMode {
		if isDir, _ := filesystem.IsDir(prj.Fs, from); isDir {
			if err := filesystem.CopyDir(prj.Fs, from, to); err != nil {
				return false, err
			}
		} else {
			if err := filesystem.CopyFile(prj.Fs, from, to); err != nil {
				return false, err
			}
		}

		io.Verbose(`<comment>Copied "%s"</comment>`, from)
	} else {
		// create symlink
		rel, _ := filepath.Rel(filepath.Dir(to), from)
		if err := prj.Fs.(afero.Symlinker).SymlinkIfPossible(rel, to); err != nil {
			return false, err
		}

		io.Verbose(`<comment>Created symlink for link "%s"</comment>`, from)
	}

	return true, nil
}
