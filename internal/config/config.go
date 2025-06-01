package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HTTPServer struct {
	Addr            string        `yaml:"address" env-required:"true"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	JWT         JWTConfig  `yaml:"jwt"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or could not load it, continuing...")
	}

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flag.StringVar(&configPath, "config", "", "path to the configuration file")
		flag.Parse()

		if configPath == "" {
			log.Fatal("Config path is not set. Please set CONFIG_PATH env var or --config flag")
		}
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	// Load YAML config into struct
	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err)
	}

	log.Printf("Loaded config from: %s", configPath)
	return &cfg
}
