package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"net/http"
	"rest_api/model"
	"testing"
)

const (
	databaseName = "movies"
	databasePort = 5432
)

var databaseHost = getEnv("DB_HOST", "localhost")
var apiHost = getEnv("API_HOST", "http://localhost:3000")

var databaseUser = getEnv("DB_USER", "user")
var databasePassword = getEnv("DB_PASS", "password")

var db *sql.DB

type AppSuite struct {
	suite.Suite
}

func (s *AppSuite) SetupSuite() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("test init failed: %s", err)
	}
	clearDatabase()
}

func (s *AppSuite) TearDownSuite() {
	defer db.Close()
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

func (s *AppSuite) TestGetAllMovies() {
	// GET when no movies exist

	response, err := http.Get(apiHost + "/movies")
	s.NoErrorf(err, "Should get no error from request initally")
	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be ok")

	responseData, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	}
	defer response.Body.Close()
	var movies []model.Movie
	jsonErr := json.Unmarshal(responseData, &movies)
	if jsonErr != nil {
		log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	}
	s.ElementsMatchf([]model.Movie{}, movies, "Should return empty list")

	// GET when movies exist in db
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})
	createMovieInDatabase(model.Movie{MovieId: "2", MovieName: "name2"})

	response, err = http.Get(apiHost + "/movies")
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be ok")

	responseData, readErr = io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	}
	defer response.Body.Close()
	movies = []model.Movie{}
	expectedMovies := []model.Movie{{MovieId: "1", MovieName: "name1"}, {MovieId: "2", MovieName: "name2"}}
	jsonErr = json.Unmarshal(responseData, &movies)
	if jsonErr != nil {
		log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	}
	s.ElementsMatchf(expectedMovies, movies, "Should return two movies")
}

func createMovieInDatabase(movie model.Movie) {
	sqlStatement := `INSERT INTO "movies" (movieId, movieName) VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, movie.MovieId, movie.MovieName)
	if err != nil {
		log.Fatalf("Failed creating movie in db for testing: %s", err)
	}
}

func clearDatabase() {
	sqlStatement := `DELETE FROM "movies";`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed creating movie in db for testing: %s", err)
	}
}
