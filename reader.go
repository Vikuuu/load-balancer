package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func readYAML(pwd string) ([]string, error) {
	configFile := filepath.Join(pwd, "load.yaml")
	f, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	server := struct {
		Data []string `yaml:"serve"`
	}{}

	err = yaml.Unmarshal(f, server)
	if err != nil {
		return nil, err
	}

	return server.Data, nil
}
