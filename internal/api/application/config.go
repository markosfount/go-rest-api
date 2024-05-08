package application

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func CreateDB() *sql.DB {
	dbHost := GetEnv("DB_HOST", "localhost")
	dbUser := GetEnv("DB_USER", "user")
	dbPassword := GetEnv("DB_PASS", "password")
	dbName := GetEnv("DB_NAME", "movies")
	// Initialise the connection pool.
	return SetupDB(dbHost, dbUser, dbPassword, dbName)
}

func SetupDB(dbHost, dbUser, dbPassword, dbName string) *sql.DB {
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		log.Fatalf("test init failed: %s", err)
	}

	return db
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
