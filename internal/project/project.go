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

func LoadProject(desc Descriptor) (prj *Project, err error) {
	if exists, _ := filesystem.DirExists(desc.Fs, desc.ProjectDir); !exists {
		return nil, errors.Errorf("Project dir \"%s\" does not exist", desc.ProjectDir)
	}

	prj = &Project{
		Fs:         desc.Fs,
		ProjectDir: desc.ProjectDir,
		ShadowDir:  filepath.Join(desc.ProjectDir, common.ShadowDir),
	}

	prj.ShadowModules, err = prj.loadShadowModules()

	if err != nil {
		return nil, err
	}

	prj.StandardModules, err = prj.loadStandardModules()

	if err != nil {
		return nil, err
	}

	return prj, nil
}

func (prj *Project) loadShadowModules() ([]*ShadowModule, error) {
	paths, err := filesystem.Glob(prj.Fs, filepath.Join(prj.ProjectDir, common.ShadowDir, "*"))

	if err != nil {
		return nil, err
	}

	var modules []*ShadowModule
	for _, path := range paths {
		cfgFilePath := filepath.Join(path, common.ShadowFile)
		if exists, _ := filesystem.Exists(prj.Fs, cfgFilePath); !exists {
			io.Verbose(`<warning>No config file found at "%s"</warning>`, cfgFilePath)
			continue
		}

		links, err := config.ReadLinks(prj.Fs, cfgFilePath)

		if err != nil {
			return nil, err
		}

		if len(links) == 0 {
			return nil, errors.Errorf(`Empty YAML file provided at "%s"`, cfgFilePath)
		}

		modules = append(modules, &ShadowModule{
			Name:      filepath.Base(path),
			ModuleDir: path,
			Links:     links,
		})
	}

	return modules, nil
}

func (prj *Project) loadStandardModules() ([]*StandardModule, error) {
	paths, err := filesystem.Glob(prj.Fs, filepath.Join(prj.ProjectDir, "*", "Pyz*", "*", "*"))

	if err != nil {
		return nil, err
	}

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
		return nil, nil
	}

	var modules []*StandardModule
	for name, directories := range desc {
		modules = append(modules, &StandardModule{
			Name:        name,
			Directories: directories,
		})
	}

	return modules, nil
}

func (prj *Project) Validate() error {
	if exists, _ := filesystem.DirExists(prj.Fs, prj.ShadowDir); !exists {
		return errors.Errorf(`Shadow directory does not exist at "%s"`, prj.ShadowDir)
	}

	return nil
}
