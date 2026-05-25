// Package config manages local DevLevel configuration persisted in the user's
// home directory at ~/.devlevel/config.json.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	dirName  = ".devlevel"
	fileName = "config.json"
)

// Config holds the persisted user configuration.
type Config struct {
	GitHubUsername string `json:"github_username"`
}

// ErrNotConfigured is returned when no configuration file exists yet.
var ErrNotConfigured = errors.New("no configuration found")

// Load reads the configuration from ~/.devlevel/config.json.
// Returns ErrNotConfigured if the file does not exist yet.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotConfigured
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to ~/.devlevel/config.json,
// creating the directory if it does not exist.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}

	return nil
}

// Dir returns the path to the ~/.devlevel directory.
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, dirName), nil
}

// configPath returns the full path to the config file.
func configPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fileName), nil
}

// ValidateUsername performs basic validation on a GitHub username.
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if len(username) > 39 {
		return errors.New("username is too long (max 39 characters)")
	}
	return nil
}
