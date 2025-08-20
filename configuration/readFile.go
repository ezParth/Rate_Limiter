package configuration

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYAML(path string) (*YAMLCONFIG, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var yamlFile YAMLCONFIG

	if err = yaml.Unmarshal(file, &yamlFile); err != nil {
		return nil, err
	}

	fmt.Printf("%v", yamlFile)
	return &yamlFile, nil
}
