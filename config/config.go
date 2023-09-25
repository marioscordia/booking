package config

import (
	"encoding/json"
	"os"
)

type DB struct {
	URI string `json:"uri"`
	Name string `json:"name"`
	User string `json:"user"`
	Hotel string `json:"hotel"`
	Room string `json:"room"`
	Book string `json:"book"`
}

type Config struct {
	Port string `json:"port"`
	Secret string `json:"secret"`
	DB DB `json:"db"`
}

func NewConfig() (*Config, error) {
	file, err := os.Open("./config/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config *Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if err := os.Setenv("JWT_SECRET", config.Secret); err != nil {
		return nil, err
	}

	return config, nil
}