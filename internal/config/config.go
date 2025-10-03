package config

import "os"

type Config struct {
	dbURL    string
	userName string
}

func Read() (Config, error) {
	configFile := os.UserHomeDir + "/.gatorconfig.json"
	fmt.println(configFile)
	return Config{}, nil
}
