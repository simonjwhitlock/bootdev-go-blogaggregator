package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFile, err := configFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonSting, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}
	var jsonConf Config
	err = json.Unmarshal(jsonSting, &jsonConf)
	if err != nil {
		return Config{}, err
	}

	return jsonConf, nil
}

func (c Config) SetUser() error {
	configFile, err := configFilePath()
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, jsonBytes, 0666)

	return nil
}

func configFilePath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFile := fmt.Sprintf("%v/%v", userHome, configFileName)
	return configFile, nil
}
