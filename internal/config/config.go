package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Config struct {
		API    API    `json:"api"`
		Client Client `json:"client"`
		Sqlite Sqlite `json:"sqlite"`
	}

	API struct {
		Host            string `json:"host"`
		Port            string `json:"port"`
		ReadTimeout     int    `json:"readTimeout"`
		WriteTimeout    int    `json:"writeTimeout"`
		MaxHeaderBytes  int    `json:"maxHeaderBytes"`
		ShutdownTimeout int    `json:"ctxTimeout"`
	}

	Client struct {
		Port string `json:"port"`
	}

	Sqlite struct {
		Driver           string `json:"driver"`
		DatabaseFileName string `json:"databaseFileName"`
		SchemePath       string `json:"schemePath"`
		ImagesPath       string `json:"imagesPath"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	var config Config

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	return &config, nil
}

func (c *Config) ServerAddress() string {
	host := c.API.Host
	port := c.API.Port
	if host == "localhost" || host == "127.0.0.1" {
		return fmt.Sprintf("%s%s", host, port)
	}
	return host
}
