package config

import (
	"github.com/andreaspenz/shadow/internal/filesystem"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Links map[string]string

func ReadLinks(fs afero.Fs, path string) (Links, error) {
	in, err := filesystem.ReadFile(fs, path)

	if err != nil {
		return nil, errors.Wrapf(err, `Unable to read file at "%s"`, path)
	}

	links := make(Links)
	err = yaml.Unmarshal(in, &links)

	if err != nil {
		return nil, errors.Wrapf(err, `Invalid YAML file provided at "%s"`, path)
	}

	return links, nil
}

func WriteLinks(fs afero.Fs, path string, links Links) error {
	out, err := yaml.Marshal(&links)

	if err != nil {
		return err
	}

	if err := filesystem.WriteFile(fs, path, out, 0666); err != nil {
		return errors.Wrapf(err, `Unable to write file at "%s"`, path)
	}

	return nil
}
