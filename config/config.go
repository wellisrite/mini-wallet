// julo/config/config.go
package config

import "os"

// Config holds configuration values.
type Config struct {
	JWTSecret string
	// Add other configuration values as needed
}

var AppConfig *Config

// Initialize loads configuration values from environment variables.
func Initialize() {
	AppConfig = &Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
		// Load other configuration values here
	}
}
