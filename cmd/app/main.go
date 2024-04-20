package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rest_api/internal/api/application"
	"rest_api/internal/api/handler"
	"rest_api/internal/api/service"
	"rest_api/internal/data"

	"github.com/gorilla/mux"
)

func main() {
	dbHost := application.GetEnv("DB_HOST", "localhost")
	dbUser := application.GetEnv("DB_USER", "user")
	dbPassword := application.GetEnv("DB_PASS", "password")
	dbName := application.GetEnv("DB_NAME", "movies")
	log.Println("ospe")
	// Initialise the connection pool.
	db := application.SetupDB(dbHost, dbUser, dbPassword, dbName)

	userRepository := data.UserRepository{DB: db}
	movieRepository := data.MovieRepository{DB: db}

	movieService := service.MovieService{MovieRepository: movieRepository}

	h := &handler.Handler{
		UserRepository: userRepository,
		MovieService:   movieService,
	}
	r := mux.NewRouter()

	r.HandleFunc("/ping", h.PingHandler).Methods(http.MethodGet)
	//r.HandleFunc("/movies", h.BasicAuth(h.GetMovies)).Methods(http.MethodGet)
	r.HandleFunc("/movies", h.GetMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies/{movieId}", h.GetMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies", h.AddMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{movieId}", h.UpdateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{movieId}", h.DeleteMovie).Methods(http.MethodDelete)
	http.Handle("/", r)

	// listen port
	err := http.ListenAndServe(":3000", nil)
	// print any server-based error messages
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
