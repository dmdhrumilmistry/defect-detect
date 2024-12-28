package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Host config
	HostPort string

	// Db Config
	DbUri          string
	DbName         string
	DbQueryTimeout int

	// App config
	AppEnv   string
	IsDevEnv bool

	// Secrets
	GithubToken string
}

var DefaultConfig = NewConfig()

func NewConfig() *Config {
	godotenv.Load()
	env := strings.ToLower(getEnvString("ENV", "production"))

	return &Config{
		HostPort:       getEnvString("HOST_PORT", "8080"),
		DbUri:          getEnvString("DB_URI", "mongodb://localhost:27017"),
		DbName:         getEnvString("DB_NAME", "defectdetect"),
		DbQueryTimeout: getEnvInt("DB_QUERY_TIMEOUT", 5),
		AppEnv:         env,
		IsDevEnv:       env == "dev",
		GithubToken:    getEnvString("GITHUB_TOKEN", ""),
	}
}

func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Error().Err(err).Msgf("failed to get env value as int: %s", key)
		return defaultValue
	}

	return value
}
