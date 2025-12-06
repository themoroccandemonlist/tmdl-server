package config

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boj/redistore"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Store       *redistore.RediStore
	Database    *pgxpool.Pool
	OAuth2      *oauth2.Config
	SessionKey  []byte
	Environment string
	Sanitizer   *bluemonday.Policy
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

func LoadPostgreSQL(ctx context.Context) *pgxpool.Pool {
	dbUser := GetConfig("DB_USER")
	dbPassword := GetConfig("DB_PASSWORD")
	dbHost := GetConfig("DB_HOST")
	dbPort := GetConfig("DB_PORT")
	dbName := GetConfig("DB_NAME")
	ssl := GetConfig("SSL")

	var sslMode string
	if ssl == "true" {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}

	dbUrl := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=%v", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal("Unable to create a connection pool: ", err)
	}
	return pool
}

func (c *Config) Close() {
	log.Printf("Cleaning up resources...\n")

	if c.Database != nil {
		c.Database.Close()
		log.Printf("PostgreSQL pool closed.\n")
	}

	if c.Store != nil {
		c.Store.Close()
		log.Printf("RediStore closed.\n")
	}

	log.Printf("All resources cleaned up.\n")
}

func New() *Config {
	gob.Register(uuid.UUID{})
	gob.Register([]string{})

	p := bluemonday.StrictPolicy()

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

	oAuth2 := &oauth2.Config{
		ClientID:     GetConfig("GOOGLE_CLIENT_ID"),
		ClientSecret: GetConfig("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		RedirectURL: GetConfig("GOOGLE_REDIRECT_URL"),
	}

	redis := LoadRedis(sessionKey, env)
	postgresql := LoadPostgreSQL(context.Background())

	return &Config{
		Store:       redis,
		Database:    postgresql,
		OAuth2:      oAuth2,
		SessionKey:  sessionKey,
		Environment: env,
		Sanitizer:   p,
	}
}
