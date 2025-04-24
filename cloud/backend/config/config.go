package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	sbytes "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/securebytes"
	sstring "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/securestring"
)

type Configuration struct {
	App   AppConfig
	Cache CacheConf
	DB    DBConfig
}

type CacheConf struct {
	URI string
}

type AppConfig struct {
	DataDirectory            string
	Port                     string
	IP                       string
	AdministrationHMACSecret *sbytes.SecureBytes
	AdministrationSecretKey  *sstring.SecureString
	GeoLiteDBPath            string
	BannedCountries          []string
	BetaAccessCode           string
}

type DBConfig struct {
	URI           string
	MapleSendName string
}

func NewProvider() *Configuration {
	var c Configuration

	// --- Application section ---
	c.App.DataDirectory = getEnv("BACKEND_APP_DATA_DIRECTORY", true)
	c.App.Port = getEnv("BACKEND_PORT", true)
	c.App.IP = getEnv("BACKEND_IP", false)
	c.App.AdministrationHMACSecret = getSecureBytesEnv("BACKEND_APP_ADMINISTRATION_HMAC_SECRET", false)
	c.App.AdministrationSecretKey = getSecureStringEnv("BACKEND_APP_ADMINISTRATION_SECRET_KEY", false)
	c.App.GeoLiteDBPath = getEnv("BACKEND_APP_GEOLITE_DB_PATH", false)
	c.App.BannedCountries = getStringsArrEnv("BACKEND_APP_BANNED_COUNTRIES", false)
	c.App.BetaAccessCode = getEnv("BACKEND_APP_BETA_ACCESS_CODE", false)

	// --- Database section ---
	c.DB.URI = getEnv("BACKEND_DB_URI", true)
	c.DB.MapleSendName = getEnv("BACKEND_DB_MAPLESEND_NAME", true)

	// --- Cache ---
	c.Cache.URI = getEnv("BACKEND_CACHE_URI", true)

	return &c
}

func getEnv(key string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return value
}

func getSecureStringEnv(key string, required bool) *sstring.SecureString {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	ss, err := sstring.NewSecureString(value)
	if ss == nil && required == false {
		return nil
	}
	if err != nil {
		log.Fatalf("Environment variable failed to secure: %v", err)
	}
	return ss
}

func getBytesEnv(key string, required bool) []byte {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return []byte(value)
}

func getSecureBytesEnv(key string, required bool) *sbytes.SecureBytes {
	value := getBytesEnv(key, required)
	sb, err := sbytes.NewSecureBytes(value)
	if sb == nil && required == false {
		return nil
	}
	if err != nil {
		log.Fatalf("Environment variable failed to secure: %v", err)
	}
	return sb
}

func getEnvBool(key string, required bool, defaultValue bool) bool {
	valueStr := getEnv(key, required)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Fatalf("Invalid boolean value for environment variable %s", key)
	}
	return value
}

func getStringsArrEnv(key string, required bool) []string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return strings.Split(value, ",")
}

func getUint64Env(key string, required bool) uint64 {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	valueUint64, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		log.Fatalf("Invalid uint64 value for environment variable %s", key)
	}
	return valueUint64
}
