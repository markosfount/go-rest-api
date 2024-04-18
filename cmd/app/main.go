package main

import (
	"fmt"
	"net/http"
	"os"
	"rest_api/internal/api/application"
	"rest_api/internal/api/handler"
	"rest_api/internal/data"

	"github.com/gorilla/mux"
)

func main() {
	dbHost := application.GetEnv("DB_HOST", "localhost")
	dbUser := application.GetEnv("DB_USER", "user")
	dbPassword := application.GetEnv("DB_PASS", "password")
	dbName := application.GetEnv("DB_NAME", "movies")

	// Initialise the connection pool.
	db := application.SetupDB(dbHost, dbUser, dbPassword, dbName)

	userRepository := data.UserRepository{DB: db}
	movieRepository := data.MovieRepository{DB: db}

	handler := &handler.Handler{
		UserRepository: userRepository,
		MovieRepository: movieRepository,
	}
	r := mux.NewRouter()

	r.HandleFunc("/ping", handler.PingHandler).Methods(http.MethodGet)
	r.HandleFunc("/movies", handler.BasicAuth(handler.GetMovies)).Methods(http.MethodGet)
	r.HandleFunc("/movies/{movieId}", handler.GetMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies", handler.AddMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{movieId}", handler.UpdateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{movieId}", handler.DeleteMovie).Methods(http.MethodDelete)
	http.Handle("/", r)

	// listen port
	err := http.ListenAndServe(":3000", nil)
	// print any server-based error messages
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
