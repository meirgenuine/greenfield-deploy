package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type APIConfig struct {
	URL string `yaml:"url"`
}

type Config struct {
	BotToken string    `yaml:"token"`
	API      APIConfig `yaml:"api"`
}

func LoadConfig() *Config {
	// Open the configuration file
	file, err := os.Open("config.yaml")
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
