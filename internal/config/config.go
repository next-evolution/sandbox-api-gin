package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBSchema   string
	DBUser     string
	DBPassword string

	RedisHost string
	RedisPort string

	JWTIssuers   []string
	JWTAudiences []string
	JWTOrigins   []string

	SessionTTL int

	ServerPort string
	GINMode    string // "debug" / "release" / "test"

	FxRateURL string

	StorageBucket        string
	StorageFX            string
	IndicatorExcludeList []string

	CsvBulkLoadSize int
	ImportCheckSkip bool
}

func Load() *Config {
	sessionTTL := 3600
	if v := os.Getenv("SESSION_TTL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			sessionTTL = n
		}
	}

	csvBulkLoadSize := 500
	if v := os.Getenv("CSV_BULK_LOAD_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			csvBulkLoadSize = n
		}
	}

	return &Config{
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "43306"),
		DBSchema:             getEnv("DB_SCHEMA", "sandbox_local"),
		DBUser:               getEnv("DB_USER", "sandbox_app"),
		DBPassword:           getEnv("DB_PASSWORD", "s4ndb0x_app"),
		RedisHost:            getEnv("REDIS_HOST", "localhost"),
		RedisPort:            getEnv("REDIS_PORT", "46379"),
		JWTIssuers:           collectEnvs("JWT_ISSUER1", "JWT_ISSUER2"),
		JWTAudiences:         collectEnvs("JWT_AUDIENCE1", "JWT_AUDIENCE2", "JWT_AUDIENCE3"),
		JWTOrigins:           collectEnvs("JWT_ORIGIN1", "JWT_ORIGIN2"),
		SessionTTL:           sessionTTL,
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		GINMode:              getEnv("GIN_MODE", "debug"),
		FxRateURL:            getEnv("FX_RATE_URL", ""),
		StorageBucket:        getEnv("STORAGE_BUCKET", "/tmp/sandbox"),
		StorageFX:            getEnv("STORAGE_FX", "fx"),
		IndicatorExcludeList: splitEnv("INDICATOR_EXCLUDE_LIST"),
		CsvBulkLoadSize:      csvBulkLoadSize,
		ImportCheckSkip:      os.Getenv("IMPORT_CHECK_SKIP") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func splitEnv(key string) []string {
	v := os.Getenv(key)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			result = append(result, s)
		}
	}
	return result
}

func collectEnvs(keys ...string) []string {
	var result []string
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			result = append(result, v)
		}
	}
	return result
}
