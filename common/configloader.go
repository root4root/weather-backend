package common

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(configPath string, c any) {

	file, fileErr := os.ReadFile(configPath)

	if fileErr != nil {
		saveConfig(configPath, c)
		fmt.Fprintf(os.Stderr, "Not found! Please edit just created file: %s\n", configPath)
		os.Exit(1)
	}

	yamlErr := yaml.Unmarshal(file, c)

	if yamlErr != nil {
		fmt.Fprintf(os.Stderr, "YAML decode error! %v\n", yamlErr)
		os.Exit(1)
	}
}

func saveConfig(configPath string, c any) {

	yamlContent, yamlErr := yaml.Marshal(c)

	if yamlErr != nil {
		fmt.Fprintf(os.Stderr, "YAML encode error! %s\n", yamlErr.Error())
		os.Exit(1)
	}

	fileErr := os.WriteFile(configPath, yamlContent, 0600)

	if fileErr != nil {
		fmt.Fprintf(os.Stderr, "YAML write to file error! %s\n", fileErr.Error())
		os.Exit(1)
	}

}
