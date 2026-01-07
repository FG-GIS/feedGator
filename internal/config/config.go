package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func Read() (Config, error) {
	confPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	confData, err := os.ReadFile(confPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(confData, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + configFileName, nil
}

func write(conf Config) error {
	confPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	// jsonBlob, err := json.Marshal(conf)
	// if err != nil {
	// 	return err
	// }
	confFile, err := os.Create(confPath)
	if err != nil {
		return err
	}
	defer confFile.Close()
	encoder := json.NewEncoder(confFile)
	encoder.SetIndent("", "  ")

	// err = encoder.Encode(jsonBlob)
	// if err != nil {
	// 	return err
	// }
	return encoder.Encode(conf)
}
