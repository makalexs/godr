package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Url string `yaml:"url"`
	} `yaml:"server"`
	DatabasePostgres struct {
		User 	 string `yaml:"user"`
		Pass 	 string `yaml:"pass"`
		Database string `yaml:"database"`
		Url 	 string `yaml:"url"`
		Port 	 string `yaml:"port"`
	} `yaml:"databasePostgres"`
	DatabaseMongo struct {
		Database string `yaml:"database"`
		Url 	 string `yaml:"url"`
		Port 	 string `yaml:"port"`
	} `yaml:"databaseMongo"`
}

func GetConfiguration() Config {
	f, err := os.Open("config.yml")
	if err != nil {
		return Config{};
	}

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{};
	}
	return cfg
}
