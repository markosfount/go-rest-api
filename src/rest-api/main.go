package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"rest_api/handler"
)

var dbHost = getEnv("DB_HOST", "localhost")
var dbUser = getEnv("DB_USER", "user")
var dbPassword = getEnv("DB_PASS", "password")
var dbName = getEnv("DB_NAME", "movies")

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
	r := mux.NewRouter()

	r.HandleFunc("/ping", env.TestHandler).Methods(http.MethodGet)
	r.HandleFunc("/movies", env.GetMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies/{movieId}", env.GetMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies", env.AddMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{movieId}", env.UpdateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{movieId}", env.DeleteMovie).Methods(http.MethodDelete)
	http.Handle("/", r)

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
