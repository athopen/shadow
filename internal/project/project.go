package project

import (
	"github.com/andreaspenz/shadow/internal/common"
	"github.com/andreaspenz/shadow/internal/config"
	"github.com/andreaspenz/shadow/internal/filesystem"
	"github.com/andreaspenz/shadow/internal/io"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"path/filepath"
)

type Descriptor struct {
	Fs         afero.Fs
	ProjectDir string
}

type Project struct {
	Fs              afero.Fs
	ProjectDir      string
	ShadowDir       string
	ShadowModules   []*ShadowModule
	StandardModules []*StandardModule
}

type ShadowModule struct {
	Name      string
	ModuleDir string
	Links     config.Links
}

type StandardModule struct {
	Name        string
	Directories []string
}

func LoadProject(desc Descriptor, fullLoad bool) (*Project, error) {
	if exists, _ := filesystem.DirExists(desc.Fs, desc.ProjectDir); !exists {
		return nil, errors.Errorf("Project dir does not exist at \"%s\"", desc.ProjectDir)
	}

	prj := &Project{
		Fs:         desc.Fs,
		ProjectDir: desc.ProjectDir,
		ShadowDir:  filepath.Join(desc.ProjectDir, common.ShadowDir),
	}

	if !fullLoad {
		return prj, nil
	}

	if exists, _ := filesystem.DirExists(prj.Fs, prj.ShadowDir); !exists {
		return nil, errors.Errorf("Shadow dir does not exist at \"%s\"", prj.ShadowDir)
	}

	if err := prj.attachShadowModules(); err != nil {
		return nil, err
	}

	if err := prj.attachStandardModules(); err != nil {
		return nil, err
	}

	return prj, nil
}

func (prj *Project) attachShadowModules() error {
	paths, _ := filesystem.Glob(prj.Fs, filepath.Join(prj.ProjectDir, common.ShadowDir, "*"))
	for _, path := range paths {
		cfgFilePath := filepath.Join(path, common.ShadowFile)
		if exists, _ := filesystem.Exists(prj.Fs, cfgFilePath); !exists {
			io.Verbose(`<warning>No config file found at "%s"</warning>`, cfgFilePath)
			continue
		}

		links, err := config.ReadLinks(prj.Fs, cfgFilePath)

		if err != nil {
			return err
		}

		if len(links) == 0 {
			return errors.Errorf(`Empty YAML file provided at "%s"`, cfgFilePath)
		}

		prj.ShadowModules = append(prj.ShadowModules, &ShadowModule{
			Name:      filepath.Base(path),
			ModuleDir: path,
			Links:     links,
		})
	}

	return nil
}

func (prj *Project) attachStandardModules() error {
	paths, _ := filesystem.Glob(prj.Fs, filepath.Join(prj.ProjectDir, "*", "Pyz*", "*", "*"))
	desc := make(map[string][]string)
	for _, path := range paths {
		// skip symlinks
		if link, _ := prj.Fs.(afero.Symlinker).ReadlinkIfPossible(path); link != "" {
			continue
		}

		// skip non directories
		if isDir, _ := filesystem.IsDir(prj.Fs, path); !isDir {
			continue
		}

		name := filepath.Base(path)
		desc[name] = append(desc[name], path)
	}

	if len(desc) == 0 {
		return nil
	}

	for name, directories := range desc {
		prj.StandardModules = append(prj.StandardModules, &StandardModule{
			Name:        name,
			Directories: directories,
		})
	}

	return nil
}
