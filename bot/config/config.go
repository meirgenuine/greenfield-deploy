package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BotToken             string              `yaml:"token"`
	DeploymentServiceURL string              `yaml:"deployment_service_url"`
	Users                map[string]struct{} `yaml:"users"`
}

func Load(path string) *Config {
	// Open the configuration file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the configuration file
	var yml map[string]interface{}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&yml)
	if err != nil {
		log.Fatal(err)
	}
	// validate fields
	v, exist := yml["token"]
	token, ok := v.(string)
	if !exist || !ok || len(token) < 1 {
		log.Fatal("invalid token")
	}
	v, exist = yml["deployment_service_url"]
	url, ok := v.(string)
	if !exist || !ok || len(url) < 1 {
		log.Fatal("invalid url")
	}

	v, exist = yml["users"]
	names, ok := v.([]interface{})
	if !exist || !ok || len(names) < 1 {
		log.Fatal("no users")
	}

	users := make(map[string]struct{})
	for _, v := range names {
		name, ok := v.(string)
		if !ok {
			log.Fatal("invalid username", v)
		}
		users[name] = struct{}{}
	}

	return &Config{
		BotToken:             token,
		DeploymentServiceURL: url,
		Users:                users,
	}
}
