package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var (
	config = new(Config)
	once   sync.Once
)

type Config struct {
	// HttpPort defines what port the server should handle incoming requests from
	HttpPort int `env:"HTTP_PORT,required,notEmpty" envDefault:"8080"`

	// DatabaseUrl ...
	DatabaseUrl string `env:"DATABASE_URL,required,notEmpty,unset"`
	Environment string `env:"ENVIRONMENT" envDefault:"test"`

	MaximumDBConn int `env:"MAX_DB_CONNECTION,required,notEmpty,unset" envDefault:"10"`
}

// LoadConfig initializes the configuration for the application and returns a pointer to a singleton configuration
func LoadConfig() *Config {
	once.Do(func() {
		if loadErr := godotenv.Load(".env"); loadErr != nil {
			log.Println("Error loading .env file - Ignore on Prod " + loadErr.Error())
		}

		// Parse environment variables to config file
		if err := env.Parse(config); err != nil {
			log.Fatalf("%+v", err)
		}
	})

	return config
}
