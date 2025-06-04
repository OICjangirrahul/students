package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Command line flags
	var command string
	var name string
	flag.StringVar(&command, "command", "", "migrate command (up/down/create)")
	flag.StringVar(&name, "name", "", "migration name (only for create command)")
	flag.Parse()

	if command == "" {
		log.Fatal("Command is required")
	}

	if command == "create" {
		if name == "" {
			log.Fatal("Name is required for create command")
		}
		createMigration(name)
		return
	}

	// Database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Initialize migrate
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	// Execute command
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration up completed")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration down completed")
	default:
		log.Fatal("Invalid command. Use up, down, or create")
	}
}

func createMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	basePath := "migrations"

	// Ensure migrations directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Fatal(err)
	}

	// Create up migration
	upFile := fmt.Sprintf("%s/%s_%s.up.sql", basePath, timestamp, name)
	if err := os.WriteFile(upFile, []byte(""), 0644); err != nil {
		log.Fatal(err)
	}

	// Create down migration
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", basePath, timestamp, name)
	if err := os.WriteFile(downFile, []byte(""), 0644); err != nil {
		log.Fatal(err)
	}

	log.Printf("Created migration files: %s, %s", upFile, downFile)
}
