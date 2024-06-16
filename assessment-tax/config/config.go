package config

import "os"

type Config struct {
	Port          string
	AdminUsername string
	AdminPassword string
	DatabaseURL   string
}

func LoadConfig() *Config {
	return &Config{
		Port:          getEnv("PORT", ""),
		AdminUsername: getEnv("ADMIN_USERNAME", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
	}

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
