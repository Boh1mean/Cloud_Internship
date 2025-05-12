package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	Port    string   `json:"port"`
	Servers []string `json:"servers"`
}

func LoadServer(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if len(cfg.Servers) == 0 {
		return nil, errors.New("no backend servers defined")
	}

	for _, addr := range cfg.Servers {
		if _, err := url.ParseRequestURI(addr); err != nil {
			return nil, fmt.Errorf("invalid server URL: %s", addr)
		}
	}
	return &cfg, nil
}
