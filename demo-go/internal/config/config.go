package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
    Address string `yaml:"address"`
}

type DatabaseConfig struct {
    DSN string `yaml:"dsn"`
}

func LoadConfig() *Config {
    file, err := os.Open("../../configs/development.yaml")
    if err != nil {
        log.Fatalf("could not open config file: %v", err)
    }
    defer file.Close()

    var cfg Config
    decoder := yaml.NewDecoder(file)
    if err := decoder.Decode(&cfg); err != nil {
        log.Fatalf("could not decode config: %v", err)
    }

    return &cfg
}
