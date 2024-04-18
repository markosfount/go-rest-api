package application

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

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