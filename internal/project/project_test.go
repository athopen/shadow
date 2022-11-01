package project

import (
	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
	"testing"
)

func TestLoadProjectWithAbsentProjectDir(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()

	result, err := LoadProject(Descriptor{Fs: fs, ProjectDir: "/some/path"}, false)

	c.Assert(result, qt.IsNil)
	c.Assert(err, qt.ErrorMatches, `Project dir does not exist at ".*"`)
}

func TestLoadProjectWithAbsentShadowDir(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/some/path", 0755)

	result, err := LoadProject(Descriptor{Fs: fs, ProjectDir: "/some/path"}, true)

	c.Assert(result, qt.IsNil)
	c.Assert(err, qt.ErrorMatches, `Shadow dir does not exist at ".*"`)
}

func TestLoadProjectWithoutModules(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/some/path/.shadow", 0755)

	result, err := LoadProject(Descriptor{Fs: fs, ProjectDir: "/some/path"}, true)

	c.Assert(result, qt.IsNotNil)
	c.Assert(err, qt.IsNil)
}

func TestLoadProjectWithEmptyShadowModuleConfigFile(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/some/path/.shadow/SomeModule", 0755)
	_, _ = fs.Create("/some/path/.shadow/SomeModule/.shadow.yml")

	result, err := LoadProject(Descriptor{Fs: fs, ProjectDir: "/some/path"}, true)

	c.Assert(result, qt.IsNil)
	c.Assert(err, qt.ErrorMatches, `Empty YAML file provided at ".*"`)
}

func TestLoadProjectWithShadowModuleConfigFile(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()

	data := `
from: to
`
	_ = afero.WriteFile(fs, "/some/path/.shadow/SomeModule/.shadow.yml", []byte(data), 0644)
	_ = fs.MkdirAll("/some/path/.shadow/SomeModule/from", 0755)

	result, err := LoadProject(Descriptor{Fs: fs, ProjectDir: "/some/path"}, true)

	c.Assert(result, qt.IsNotNil)
	c.Assert(result.ShadowModules, qt.HasLen, 1)
	c.Assert(result.ShadowModules[0].Name, qt.Equals, "SomeModule")
	c.Assert(result.ShadowModules[0].ModuleDir, qt.Equals, "/some/path/.shadow/SomeModule")
	c.Assert(result.ShadowModules[0].Links, qt.HasLen, 1)
	c.Assert(err, qt.IsNil)
}
