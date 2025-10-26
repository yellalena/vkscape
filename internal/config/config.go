package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AuthMethod string

const (
	AuthMethodAppToken AuthMethod = "app_token"
	AuthMethodUser     AuthMethod = "user"
)

type AuthConfig struct {
	AuthMethod   AuthMethod `json:"auth_method"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token,omitempty"`
}

func SaveConfig(cfg *AuthConfig) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadConfig() (*AuthConfig, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AuthConfig
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(configDir, "vkscape")

	err = os.MkdirAll(appDir, 0750)
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, "config.json"), nil
}
