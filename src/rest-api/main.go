package main // Update the AllBooks function so it accepts the connection pool as a

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rest_api/handler"
)

var dbHost = getEnv("DB_HOST", "localhost")
var dbUser = getEnv("DB_USER", "user")
var dbPassword = getEnv("DB_PASS", "password")
var dbName = getEnv("DB_NAME", "postgres")

// DB set up
func setupDB() *sql.DB {
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		log.Fatalf("test init failed: %s", err)
	}

	return db
}

func main() {
	// Initialise the connection pool.
	db := setupDB()

	// Create an instance of Env containing the connection pool.
	env := &handler.Env{Db: db}

	http.HandleFunc("/movies", env.GetMovies)
	http.HandleFunc("/movie", env.GetMovie)
	http.HandleFunc("/movie/create", env.AddMovie)
	http.HandleFunc("/movie/update", env.UpdateMovie)
	http.HandleFunc("/ping", env.TestHandler)

	// listen port
	err := http.ListenAndServe(":3000", nil)
	// print any server-based error messages
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
