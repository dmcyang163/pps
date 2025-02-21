package config

import (
	"encoding/json"
	"os"
)

// Config defines the node configuration.
type Config struct {
	Port         int      `json:"port"`
	SeedNodes    []string `json:"seed_nodes"`
	MaxPeers     int      `json:"max_peers"`
	PingInterval int      `json:"ping_interval"`
	DataDir      string   `json:"data_dir"` // Directory for storing file transfer data
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &Config{}
	err = decoder.Decode(cfg)
	return cfg, err
}
