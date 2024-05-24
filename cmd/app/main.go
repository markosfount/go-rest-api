package main

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest_api/internal/api/application"
	"rest_api/internal/api/config"
	"rest_api/internal/api/handler"
	"rest_api/internal/api/kafka"
	"rest_api/internal/api/service"
	"rest_api/internal/api/tmdb"
	"rest_api/internal/data"
	"rest_api/internal/scheduler"
	"sync"
	"syscall"
	"time"
)

func main() {
	// create kafka topic
	admin, err := sarama.NewClusterAdmin([]string{config.BrokerLink}, sarama.NewConfig())
	if err != nil {
		log.Fatalf("Could not create kafka admin: %v", err)
	}
	err = admin.CreateTopic(config.Topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	var tErr *sarama.TopicError
	if err != nil {
		if !errors.As(err, &tErr) || !errors.Is(tErr.Unwrap(), sarama.ErrTopicAlreadyExists) {
			log.Fatalf("Could not create topic: %v", err)
		}
	}
	// create kafka publisher
	var publisher kafka.Publisher
	if config.SyncPublish {
		publisher = &kafka.SyncPublisher{}
	} else {
		publisher = &kafka.AsyncPublisher{}
	}
	publisher.Configure(config.Topic)

	//create db and service layer
	db := application.CreateDB()

	userRepository := &data.UserRepository{DB: db}
	movieRepository := &data.MovieRepository{DB: db}

	movieService := service.NewMovieService(movieRepository)
	tmdbService := tmdb.NewService(tmdb.ApiUrl)

	r := mux.NewRouter()

	h := &handler.Handler{
		UserRepository: userRepository,
		MovieService:   movieService,
		TmdbService:    tmdbService,
		Publisher:      publisher,
	}

	r.HandleFunc("/ping", h.PingHandler).Methods(http.MethodGet)
	//r.HandleFunc("/movies", h.BasicAuth(h.GetMovies)).Methods(http.MethodGet)
	r.HandleFunc("/movies", h.GetMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies/{movieId}", h.GetMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies", h.AddMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{movieId}", h.UpdateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{movieId}", h.DeleteMovie).Methods(http.MethodDelete)

	server := http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan bool)
	defer close(done)
	quit := make(chan os.Signal, 1)

	wg := sync.WaitGroup{}

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
			// FIXME
			log.Println("Second signal caught. Shutting down NOW")
		case <-delay.C:
		}
		log.Printf("shutting down server")
		// termination delay in both context and signal listening?
		ctx, cancel := context.WithTimeout(context.Background(), terminationDelay)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("could not shutdown server: %v", err)
		}
		done <- true
		close(quit)
	}()
	sch := scheduler.NewScheduler(5, done, &wg)
	// run scheduler in background
	wg.Add(1)
	go func() {
		sch.Run()
	}()

	listenAddr := ":3000"
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}
	log.Printf("Started server on %s", listenAddr)
	wg.Wait()

	log.Println("Server stopped")

}
