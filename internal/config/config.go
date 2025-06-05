package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type HTTPServer struct {
	Addr    string        `yaml:"addr" env:"HTTP_ADDR"`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret" env:"JWT_SECRET"`
	Expiration string `yaml:"expiration" env:"JWT_EXPIRATION"`
}

type AWSConfig struct {
	Region      string `yaml:"region" env:"AWS_REGION"`
	S3Bucket    string `yaml:"s3_bucket" env:"AWS_S3_BUCKET"`
	DynamoTable string `yaml:"dynamo_table" env:"AWS_DYNAMO_TABLE"`
}

type LogConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL"`
	Format string `yaml:"format" env:"LOG_FORMAT"`
}

type Config struct {
	Env        string     `yaml:"env" env:"ENV"`
	DB         DBConfig   `yaml:"db"`
	HTTPServer HTTPServer `yaml:"http_server"`
	JWT        JWTConfig  `yaml:"jwt"`
	AWS        AWSConfig  `yaml:"aws"`
	Log        LogConfig  `yaml:"log"`
}

// LoadConfig loads configuration from both YAML file and environment variables.
// Environment variables take precedence over YAML configuration.
func LoadConfig(path string) (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't return error if .env doesn't exist
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Initialize default configuration
	cfg := &Config{
		Env: getEnv("ENV", "dev"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "students"),
		},
		HTTPServer: HTTPServer{
			Addr:    getEnv("HTTP_ADDR", ":8082"),
			Timeout: time.Duration(getEnvAsInt("HTTP_TIMEOUT", 30)) * time.Second,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-here"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		AWS: AWSConfig{
			Region:      getEnv("AWS_REGION", "us-west-2"),
			S3Bucket:    getEnv("AWS_S3_BUCKET", ""),
			DynamoTable: getEnv("AWS_DYNAMO_TABLE", ""),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// If YAML path is provided, load and merge YAML configuration
	if path != "" {
		if err := loadYAMLConfig(path, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// loadYAMLConfig loads configuration from YAML file and merges it with existing config
func loadYAMLConfig(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil
}

func LoadAWSConfig(cfg *Config) (aws.Config, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return awsCfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := time.ParseDuration(value); err == nil {
			return int(intValue.Seconds())
		}
	}
	return defaultValue
}
