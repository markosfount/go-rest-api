package main

import (
	"fmt"
	"net/http"
	"os"
	"rest_api/internal/api/application"
	"rest_api/internal/api/handler"

	"github.com/gorilla/mux"
)

func main() {
	dbHost := application.GetEnv("DB_HOST", "localhost")
	dbUser := application.GetEnv("DB_USER", "user")
	dbPassword := application.GetEnv("DB_PASS", "password")
	dbName := application.GetEnv("DB_NAME", "movies")

	// Initialise the connection pool.
	db := application.SetupDB(dbHost, dbUser, dbPassword, dbName)

	// Create an instance of Env containing the connection pool.
	env := &handler.Env{Db: db}
	r := mux.NewRouter()

	r.HandleFunc("/ping", env.TestHandler).Methods(http.MethodGet)
	r.HandleFunc("/movies", env.BasicAuth(env.GetMovies)).Methods(http.MethodGet)
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
