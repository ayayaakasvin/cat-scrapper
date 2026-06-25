package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	// "github.com/joho/godotenv"
)

const (
	configPathEnvKey = "CONFIG_PATH"
)

// Config represents the configuration structure
type Config struct {
	// Add your configuration fields here
	// Example:
	// Port int `yaml:"port"`
	HTTPServerConfig `yaml:"http-server"`
	CorsConfig       `yaml:"cors"`
	Logger           LoggerConfig `yaml:"logger"`
	SavePath         string       `yaml:"save_path"`
	SqLiteConfig     SQLiteConfig `yaml:"sqlite"`
}

type LoggerConfig struct {
	Env     string `yaml:"env" env-required:"true"`
	Service string `yaml:"service" env-required:"true"`
	// 	const (
	// 	LevelDebug Level = -4
	// 	LevelInfo  Level = 0
	// 	LevelWarn  Level = 4
	// 	LevelError Level = 8
	// )
	Level int  `yaml:"level" env-default:"0"`
	JSON  bool `yaml:"json" env-default:"false"`
}

type CorsConfig struct {
	AllowedOrigins     []string `yaml:"allowed_origins"`
	AllowedMethods     []string `yaml:"allowed_methods"`
	AllowedHeaders     []string `yaml:"allowed_headers"`
	AllowedCredentials bool     `yaml:"allow_credentials"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type SQLiteConfig struct {
	FilePath string `yaml:"db_name"`
}

// MustLoadConfig loads the configuration from the specified path
func MustLoadConfig() *Config {
	// if err := godotenv.Load(); err != nil {
	// 	log.Println(".env not found, falling back to local variables")
	// }

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
