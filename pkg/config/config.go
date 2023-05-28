package config

import (
	"fmt"
	"os"
)

type Config struct {
	SlackSigningSecret string
	AppApiEndpoint     string
	AppApiApiKey       string
}

func New() (Config, error) {
	cfg := Config{
		SlackSigningSecret: os.Getenv("SLACK_SIGNING_SECRET"),
		AppApiEndpoint:     os.Getenv("APP_API_ENDPOINT"),
		AppApiApiKey:       os.Getenv("APP_API_API_KEY"),
	}

	if cfg.SlackSigningSecret == "" {
		return Config{}, fmt.Errorf("SLACK_SIGNING_SECRET is empty")
	}

	if cfg.AppApiEndpoint == "" {
		return Config{}, fmt.Errorf("APP_API_ENDPOINT is empty")
	}

	if cfg.AppApiApiKey == "" {
		return Config{}, fmt.Errorf("APP_API_API_KEY is empty")
	}

	return cfg, nil
}
