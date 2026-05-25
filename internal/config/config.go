package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	configPathEnvKey  = "CONFIG_PATH"
	postgresURLEnvKey = "POSTGRES_URL"
	valkeyURLEnvKey   = "VALKEY_URL"
	jwtSecretKey      = "JWT_SECRET"
	gateawaySecretKey = "X_GATEAWAY_KEY"
)

// Config represents the configuration structure
type Config struct {
	// Add your configuration fields here
	// Example:
	// Port int `yaml:"port"`
	HTTPServer `yaml:"http-server"`

	CorsConfig `yaml:"cors"`

	GateawaySecret string
}

type CorsConfig struct {
	AllowedOrigins     []string `yaml:"allowed_origins"`
	AllowedMethods     []string `yaml:"allowed_methods"`
	AllowedHeaders     []string `yaml:"allowed_headers"`
	AllowedCredentials bool     `yaml:"allow_credentials"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"						env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" 															env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" 														env-required:"true"`
}

// MustLoadConfig loads the configuration from the specified path
func MustLoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, falling back to local variables")
	}

	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	return &cfg
}