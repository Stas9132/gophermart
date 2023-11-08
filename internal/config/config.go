package config

import (
	"flag"
	"os"
)

type Config struct {
	Host        string
	Port        string
	DatabaseURI string
}

func New() *Config {
	config := &Config{
		Host:        getEnv("HOST", "localhost"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURI: getEnv("DATABASE_URI", "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"),
	}
	flag.StringVar(&config.Host, "host", getEnv("HOST", "localhost"), "Address of the HTTP server")
	flag.StringVar(&config.Port, "port", getEnv("PORT", "8080"), "Listening port number")
	flag.StringVar(&config.DatabaseURI, "database-uri", getEnv("DATABASE_URI", "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"), "Database URI")

	flag.Parse()
	return config
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
