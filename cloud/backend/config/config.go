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
	App               AppConfig
	Cache             CacheConf
	DB                DBConfig
	AWS               AWSConfig
	PAPERCLOUDMailgun MailgunConfig
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
	URI            string
	MapleAuthName  string
	VaultName      string
	PaperCloudName string
}

type MailgunConfig struct {
	APIKey           string
	Domain           string
	APIBase          string
	SenderEmail      string
	MaintenanceEmail string
	FrontendDomain   string
	BackendDomain    string
}

type AWSConfig struct {
	AccessKey  string
	SecretKey  string
	Endpoint   string
	Region     string
	BucketName string
}

func NewProvider() *Configuration {
	var c Configuration

	// --------- SHARED ------------
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
	c.DB.MapleAuthName = getEnv("BACKEND_DB_MAPLEAUTH_NAME", true)
	c.DB.VaultName = getEnv("BACKEND_DB_VAULT_NAME", true)
	c.DB.PaperCloudName = getEnv("BACKEND_DB_PAPERCLOUD_NAME_NAME", true)

	// --- Cache ---
	c.Cache.URI = getEnv("BACKEND_CACHE_URI", true)

	// --- AWS ---
	c.AWS.AccessKey = getEnv("BACKEND_AWS_ACCESS_KEY", true)
	c.AWS.SecretKey = getEnv("BACKEND_AWS_SECRET_KEY", true)
	c.AWS.Endpoint = getEnv("BACKEND_AWS_ENDPOINT", true)
	c.AWS.Region = getEnv("BACKEND_AWS_REGION", true)
	c.AWS.BucketName = getEnv("BACKEND_AWS_BUCKET_NAME", true)

	// --------- PaperCloud ------------
	// --- Mailgun ---
	c.PAPERCLOUDMailgun.APIKey = getEnv("BACKEND_PAPERCLOUD_MAILGUN_API_KEY", true)
	c.PAPERCLOUDMailgun.Domain = getEnv("BACKEND_PAPERCLOUD_MAILGUN_DOMAIN", true)
	c.PAPERCLOUDMailgun.APIBase = getEnv("BACKEND_PAPERCLOUD_MAILGUN_API_BASE", true)
	c.PAPERCLOUDMailgun.SenderEmail = getEnv("BACKEND_PAPERCLOUD_MAILGUN_SENDER_EMAIL", true)
	c.PAPERCLOUDMailgun.MaintenanceEmail = getEnv("BACKEND_PAPERCLOUD_MAILGUN_MAINTENANCE_EMAIL", true)
	c.PAPERCLOUDMailgun.FrontendDomain = getEnv("BACKEND_PAPERCLOUD_MAILGUN_FRONTEND_DOMAIN", true)
	c.PAPERCLOUDMailgun.BackendDomain = getEnv("BACKEND_PAPERCLOUD_MAILGUN_BACKEND_DOMAIN", true)

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
