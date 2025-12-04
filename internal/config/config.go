package config

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

type Config struct {
	Store *redistore.RediStore
}

const (
	Production  = "PRODUCTION"
	Development = "DEVELOPMENT"
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

func LoadRedis(sessionKey []byte, env string) *redistore.RediStore {
	addr := GetConfig("REDIS_ADDRESS")
	pass := GetConfig("REDIS_PASSWORD")
	poolSize, _ := strconv.Atoi(GetConfig("REDIS_POOL_SIZE"))

	store, err := redistore.NewRediStore(poolSize, "tcp", addr, "", pass, sessionKey)
	if err != nil {
		log.Fatal("Unable to create a RediStore: ", err)
	}

	secureCookies := env == Production

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   secureCookies,
		SameSite: http.SameSiteLaxMode,
	}

	return store
}

func New() *Config {

	var env string
	err := godotenv.Load()
	if err != nil {
		log.Printf("Couldn't load .env, we must be in production.\n")
		env = Production
	} else {
		env = Development
	}

	sessionKeyStr := GetConfig("SESSION_KEY")
	var sessionKey []byte
	if sessionKeyStr != "" {
		sessionKey = []byte(sessionKeyStr)
	}

	redis := LoadRedis(sessionKey, env)

	return &Config{
		Store: redis,
	}
}
