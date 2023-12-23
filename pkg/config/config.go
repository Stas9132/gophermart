package config

import (
	"flag"
	"os"
)

type Config struct {
	Address              string
	DatabaseURI          string
	AccuralSystemAddress string
}

func New() *Config {
	config := &Config{
		Address:              getEnv("RUN_ADDRESS", ":8080"),
		DatabaseURI:          getEnv("DATABASE_URI", "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"),
		AccuralSystemAddress: getEnv("ACCRUAL_SYSTEM_ADDRESS", ""),
	}
	flag.StringVar(&config.Address, "a", getEnv("RUN_ADDRESS", ":8080"), "Address of the HTTP server")
	flag.StringVar(&config.DatabaseURI, "d", getEnv("DATABASE_URI", "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"), "Database URI")

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
