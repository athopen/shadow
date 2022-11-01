package cmd

import (
	"github.com/andreaspenz/shadow/internal/filesystem"
	"github.com/andreaspenz/shadow/internal/io"
	"github.com/andreaspenz/shadow/internal/project"
	"github.com/symfony-cli/console"
	"io/fs"
	"os"
	"path/filepath"
)

func newCleanCommand() *console.Command {
	return &console.Command{
		Name:   "clean",
		Usage:  "Remove broken links",
		Action: adoptActionFunc(cleanAction, true),
	}
}

func cleanAction(_ *console.Context, prj *project.Project) error {
	var removed int
	err := filesystem.Walk(prj.Fs, prj.ProjectDir, func(path string, info fs.FileInfo, err error) error {
		// check if file is symlink
		if info.Mode()&os.ModeSymlink != os.ModeSymlink {
			return nil
		}

		// check if symlink is broken
		if _, err := filepath.EvalSymlinks(path); err == nil {
			return nil
		}

		// remove broken symlink
		if err := prj.Fs.Remove(path); err != nil {
			return err
		}

		removed++

		io.Verbose(`<comment>Removed broken link "%s"</comment>`, path)

		return nil
	})

	if err != nil {
		return nil
	}

	if removed == 0 {
		io.Write("<warning>No broken links found</warning>")
		return nil
	}

	io.Write(`<info>Removed %d broken link(s)</info>`, removed)

	return nil
}
