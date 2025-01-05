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
	AppEnv              string
	IsDevEnv            bool
	DefaultWorkersCount int

	// Secrets
	GithubToken string

	// Analyzer Config
	RunOsv  bool
	RunMpaf bool
	RunEpss bool
}

var DefaultConfig = NewConfig()

func NewConfig() *Config {
	godotenv.Load()
	env := strings.ToLower(getEnvString("ENV", "production"))

	return &Config{
		HostPort:            getEnvString("HOST_PORT", "8080"),
		DbUri:               getEnvString("DB_URI", "mongodb://localhost:27017"),
		DbName:              getEnvString("DB_NAME", "defectdetect"),
		DbQueryTimeout:      getEnvInt("DB_QUERY_TIMEOUT", 5),
		AppEnv:              env,
		IsDevEnv:            env == "dev",
		DefaultWorkersCount: getEnvInt("DEFAULT_WORKERS_COUNT", 30),
		GithubToken:         getEnvString("GITHUB_TOKEN", ""),

		// Analyzer Config
		RunOsv:  getEnvBool("RUN_OSV_ANALYZER"),
		RunMpaf: getEnvBool("RUN_MPAF_ANALYZER"),
		RunEpss: getEnvBool("RUN_EPSS_ANALYZER"),
	}
}

func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvBool(key string) bool {
	value := strings.ToLower(os.Getenv(key))
	return value == "true" || value == "1" || value == "yes"
}

func getEnvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Error().Err(err).Msgf("failed to get env value as int: %s", key)
		return defaultValue
	}

	return value
}
