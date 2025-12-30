package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ComiketBackendConfigFileConfig struct {
	LogFilePath string `yaml:"logFilePath"`
}

type ComiketBackendConfig struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
	Logging struct {
		LogLevel string                         `yaml:"logLevel"`
		File     ComiketBackendConfigFileConfig `yaml:"file"`
	} `yaml:"logging"`
	Db struct {
		Postgres struct {
			Host         string `yaml:"host"`
			Port         int    `yaml:"port"`
			DatabaseName string `yaml:"databaseName"`
			Username     string `yaml:"username"`
			Password     string `yaml:"password"`
		} `yaml:"postgres"`
	} `yaml:"db"`
}

func LoadConfigFromFile(path string) (*ComiketBackendConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var comiketBackendConfig ComiketBackendConfig
	yamlDecoder := yaml.NewDecoder(f)
	err = yamlDecoder.Decode(&comiketBackendConfig)
	if err != nil {
		return nil, err
	}

	return &comiketBackendConfig, nil
}
