package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ReadSourceYaml(path string, dest *SourceConfig) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(dest); err != nil {
		return errors.WithStack(err)
	}

	logger.Info("success to read source file", "path", path)
	return nil
}

func ReadCacheJson(path string, dest *CacheConfig) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(dest); err != nil {
		return errors.WithStack(err)
	}

	logger.Info("success to read cache file", "path", path)
	return nil
}

func WriteCacheJson(path string, src CacheConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return errors.WithStack(err)
	}
	file, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(src); err != nil {
		return errors.WithStack(err)
	}

	logger.Info("success to write cache file", "path", path)
	return nil
}

func WriteDestMarkdown(path string, artifact io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return errors.WithStack(err)
	}
	file, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if _, err = io.Copy(file, artifact); err != nil {
		return errors.WithStack(err)
	}

	logger.Info("success to write markdown file", "path", path)
	return nil
}
