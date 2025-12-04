package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetConfig(key string) string {
	value, err := ReadSecret(key)
	if err == nil {
		log.Printf("Loaded '%v' from Docker secrets.\n", key)
		return value
	}

	value = os.Getenv(key)
	if value != "" {
		log.Printf("Loaded '%v' from environment variables.\n", key)
		return value
	}

	return ""
}

func ReadSecret(secret string) (string, error) {
	file := filepath.Join("/run/secrets/", secret)
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
