package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Retrieve environment variables
    dbHost := os.Getenv("DB_HOST")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbName := os.Getenv("DB_NAME")

    // Construct connection string
    connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbName, dbPassword)

    // Open database connection
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database connection: %v", err)
    }

    // Ping database to verify connection
    if err = DB.Ping(); err != nil {
        log.Fatalf("Error pinging database: %v", err)
    }

    fmt.Println("Successfully connected to the database!")
}