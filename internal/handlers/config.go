package config

import "os"

type Config struct {
	MattermostURL   string
	MattermostToken string
	TarantoolAddr   string
}

func Load() (*Config, error) {
	return &Config{
		MattermostURL:   getEnv("MATTERMOST_URL", "http://localhost:8065"),
		MattermostToken: getEnv("MATTERMOST_TOKEN", ""),
		TarantoolAddr:   getEnv("TARANTOOL_ADDR", "localhost:3301"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
