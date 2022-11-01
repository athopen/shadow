package cmd

import (
	"github.com/andreaspenz/shadow/internal/common"
	"github.com/andreaspenz/shadow/internal/config"
	"github.com/andreaspenz/shadow/internal/filesystem"
	"github.com/andreaspenz/shadow/internal/io"
	"github.com/andreaspenz/shadow/internal/project"
	"github.com/pkg/errors"
	"github.com/symfony-cli/console"
	"os"
	"path/filepath"
	"strings"
)

type reverseCommand struct {
	cc *console.Command
}

func newReverseCommand() *reverseCommand {
	cmd := &reverseCommand{}

	cmd.cc = &console.Command{
		Name:   "reverse",
		Usage:  "Shadow Pyz modules",
		Action: adoptActionFunc(cmd.reverseAction, true),
	}

	return cmd
}

func (r *reverseCommand) reverseAction(_ *console.Context, prj *project.Project) error {
	if len(prj.StandardModules) == 0 {
		return errors.New("No valid modules found for reversion")
	}

	var reversed int
	for _, module := range prj.StandardModules {
		proceeded, err := r.reverseModule(prj, module)

		if err != nil {
			return err
		}

		if proceeded {
			reversed++
		}
	}

	if reversed == 0 {
		io.Write("<warning>No new modules found for revision</warning>")
		return nil
	}

	io.Write(`<info>Successfully reversed "%d" modules(s)</info>`, reversed)

	return nil
}

func (r *reverseCommand) reverseModule(prj *project.Project, module *project.StandardModule) (bool, error) {
	var reversed bool

	for _, path := range module.Directories {
		proceeded, err := r.reverseDirectory(prj, path)

		if err != nil {
			return false, err
		}

		if proceeded {
			reversed = true
		}
	}

	return reversed, nil
}

func (r *reverseCommand) reverseDirectory(prj *project.Project, path string) (bool, error) {
	base, namespace, layer, name := r.getInfoFromPath(prj, path)
	shadowPath := filepath.Join(prj.ShadowDir, name, base, namespace, layer, name)

	// skip if dest path already exists
	if exists, _ := filesystem.DirExists(prj.Fs, shadowPath); exists {
		io.Verbose(
			`<warning>Module "%s" with layer "%s" in "%s" already exists</warning>`,
			namespace,
			layer,
			base,
		)

		return false, nil
	}

	// create parent directories
	if err := prj.Fs.MkdirAll(filepath.Dir(shadowPath), 0755); err != nil {
		return false, errors.Errorf(`Unable to create directory at "%s"`, filepath.Dir(shadowPath))
	}

	// move to shadow folder
	if err := prj.Fs.Rename(path, shadowPath); err != nil {
		return false, errors.Errorf(`Unable to move directory "%s" to "%s"`, path, shadowPath)
	}

	shadowFile := filepath.Join(prj.ShadowDir, name, common.ShadowFile)

	// create the config file
	if _, err := prj.Fs.OpenFile(shadowFile, os.O_CREATE, 0666); err != nil {
		return false, errors.Wrapf(err, `Unable to create file at "%s"`, shadowFile)
	}

	// read config file
	links, err := config.ReadLinks(prj.Fs, shadowFile)

	if err != nil {
		return false, err
	}

	link, err := filepath.Rel(prj.ProjectDir, path)

	if err != nil {
		return false, err
	}

	links[link] = link

	if err := config.WriteLinks(prj.Fs, shadowFile, links); err != nil {
		return false, err
	}

	return true, nil
}

func (r *reverseCommand) getInfoFromPath(prj *project.Project, path string) (string, string, string, string) {
	rel, _ := filepath.Rel(prj.ProjectDir, path)
	parts := strings.Split(rel, string(filepath.Separator))

	return parts[0], parts[1], parts[2], parts[3]
}
