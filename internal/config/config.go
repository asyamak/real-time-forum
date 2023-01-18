package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Config struct {
		API      API      `json:"api"`
		Client   Client   `json:"client"`
		Database Database `json:"database"`
	}

	API struct {
		Host           string `json:"host"`
		Port           string `json:"port"`
		ReadTimeout    int    `json:"readTimeout"`
		WriteTimeout   int    `json:"writeTimeout"`
		MaxHeaderBytes int    `json:"maxHeaderBytes"`
	}

	Client struct {
		Port string `json:"port"`
	}

	Database struct {
		Driver       string `json:"driver"`
		DatabaseName string `json:"databaseName"`
		SchemePath   string `json:"schemePath"`
		ImagesPath   string `json:"imagesPath"`
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
		return fmt.Sprintf("%s:%s", host, port)
	}
	return host
}
