package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"log"
	"os"

	"github.com/joho/godotenv"

	"fmt"
)

func setupDatabase() *sql.DB {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection parameters from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	// Create a connection string
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPassword)

	// Open a database connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	return db
}
