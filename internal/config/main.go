package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilepath()
	if err != nil {
		return Config{}, err
	}

	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

const configFileName = ".gatorconfig.json"

func getConfigFilepath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, configFileName), nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilepath()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, bytes, 0644)
	return err
}
