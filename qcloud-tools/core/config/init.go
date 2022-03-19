package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func getConfigFile() (string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	file := fmt.Sprintf("%s/config/config.yaml", rootPath)
	return file, nil
}

func init() {
	file, err := getConfigFile()
	if err != nil {
		fmt.Printf("failed get config file : %v\n", err)
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("failed to read yaml file : %v\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(content, &QcloudTool)
	if err != nil {
		fmt.Printf("failed to unmarshal yaml file : %v\n", err)
	}
}
