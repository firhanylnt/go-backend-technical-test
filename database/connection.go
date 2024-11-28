package database

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

var DB *sql.DB

func Connect() {
    var err error

    err = godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    username := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbname := os.Getenv("DB_NAME")

    dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname

    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatalf("Database unreachable: %v", err)
    }

    log.Println("Connected to database!")
}
