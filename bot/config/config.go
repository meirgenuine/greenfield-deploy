package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type APIConfig struct {
	URL string `yaml:"url"`
}

type UsersConfig struct {
	Allowed map[string]string `yaml:"allowed"`
}

type Config struct {
	BotToken string      `yaml:"token"`
	API      APIConfig   `yaml:"api"`
	Users    UsersConfig `yaml:"users"`
}

func Load(path string) *Config {
	// Open the configuration file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the configuration file
	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
