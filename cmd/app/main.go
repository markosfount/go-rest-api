package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest_api/internal/api/application"
	"rest_api/internal/api/handler"
	"rest_api/internal/api/service"
	"rest_api/internal/data"
	"syscall"
	"time"
)

func main() {
	db := application.CreateDB()

	userRepository := data.UserRepository{DB: db}
	movieRepository := data.MovieRepository{DB: db}

	movieService := service.MovieService{MovieRepository: &movieRepository}

	r := mux.NewRouter()

	h := &handler.Handler{
		UserRepository: userRepository,
		MovieService:   movieService,
	}

	r.HandleFunc("/ping", h.PingHandler).Methods(http.MethodGet)
	//r.HandleFunc("/movies", h.BasicAuth(h.GetMovies)).Methods(http.MethodGet)
	r.HandleFunc("/movies", h.GetMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies/{movieId}", h.GetMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies", h.AddMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{movieId}", h.UpdateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{movieId}", h.DeleteMovie).Methods(http.MethodDelete)

	server := http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	terminationDelay := 5 * time.Second

	go func() {
		sig := <-quit
		server.SetKeepAlivesEnabled(false)
		log.Printf("Signal %v caught. Shutting down in %vs", sig, terminationDelay)

		delay := time.NewTicker(terminationDelay)
		defer delay.Stop()
		select {
		case <-quit:
			log.Println("Second signal caught. Shutting down NOW")
		case <-delay.C:
		}

		ctx, cancel := context.WithTimeout(context.Background(), terminationDelay)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("could not shutdown server: %v", err)
		}
		close(done)
	}()

	listenAddr := ":3000"
	log.Printf("Started server on %s", listenAddr)
	var err error
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	log.Println("Server stopped")

}
