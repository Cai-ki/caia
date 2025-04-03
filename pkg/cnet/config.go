package cnet

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Ip             string `json:"ip"`
	Port           int    `json:"port"`
	MaxConnections int    `json:"max_connections"`
	// ListenTimeoutSec  int    `json:"listen_timeout_sec"`
	// ReadTimeoutSec    int    `json:"read_timeout_sec"`
	// WriteTimeoutSec   int    `json:"write_timeout_sec"`
	ListenDeadlineMs int `json:"listen_deadline_ms"`
	ReadDeadlineMs   int `json:"read_deadline_ms"`
	WriteDeadlineMs  int `json:"write_deadline_ms"`
	TLS              struct {
		Enabled  bool   `json:"enabled"`
		CertFile string `json:"cert_file"`
		KeyFile  string `json:"key_file"`
	} `json:"tls"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("invalid port number")
	}
	if c.TLS.Enabled && (c.TLS.CertFile == "" || c.TLS.KeyFile == "") {
		return errors.New("TLS requires a certificate file")
	}
	return nil
}
