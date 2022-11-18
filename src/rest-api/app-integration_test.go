package main

import (
	"bytes"
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
	clearDatabase()
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

	clearDatabase()
}

func (s *AppSuite) TestGetMovie() {
	// GET when movie does not exist

	response, err := http.Get(apiHost + "/movies/1")
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusNotFound, response.StatusCode, "Expected status to be not found")

	responseData, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	}
	defer response.Body.Close()
	responseMessage := model.ResponseMessage{}
	expectedMessage := model.ResponseMessage{Message: "No movie with provided id exists"}
	jsonErr := json.Unmarshal(responseData, &responseMessage)
	if jsonErr != nil {
		log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for not found")

	// GET when movies exist in db
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})

	response, err = http.Get(apiHost + "/movies/1")
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be ok")

	responseData, readErr = io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	}
	defer response.Body.Close()
	movie := model.Movie{}
	expectedMovie := model.Movie{MovieId: "1", MovieName: "name1"}
	jsonErr = json.Unmarshal(responseData, &movie)
	if jsonErr != nil {
		log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	}
	s.Equal(expectedMovie, movie, "Should return created movie")

	clearDatabase()
}

func (s *AppSuite) TestCreateMovie() {
	// Create movie
	movieId := "1"
	movieToCreate := model.Movie{MovieId: movieId, MovieName: "name1"}
	body, jsonErr := json.Marshal(movieToCreate)

	response, err := http.Post(apiHost+"/movies", "application/json", bytes.NewBuffer(body))
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusCreated, response.StatusCode, "Expected status to created")

	responseData, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	}
	defer response.Body.Close()
	movie := model.Movie{}
	expectedMovie := model.Movie{MovieId: "1", MovieName: "name1"}
	jsonErr = json.Unmarshal(responseData, &movie)
	if jsonErr != nil {
		log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	}
	s.Equal(expectedMovie, movie, "Should return created movie")
	// check that movie was created in db
	savedMovie := getMovieFromDatabase(movieId)
	s.Equal(expectedMovie, savedMovie, "Should return created movie")

	// Try to create already existing movie
	clearDatabase()
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})

	response, err = http.Post(apiHost+"/movies", "application/json", bytes.NewBuffer(body))
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusConflict, response.StatusCode, "Expected status to be conflict")

	responseData, readErr = io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	}
	defer response.Body.Close()
	responseMessage := model.ResponseMessage{}
	expectedMessage := model.ResponseMessage{Message: "A movie with the provided id already exists"}
	jsonErr = json.Unmarshal(responseData, &responseMessage)
	if jsonErr != nil {
		log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for conflict")

	clearDatabase()
}

func (s *AppSuite) TestUpdateMovie() {
	// Update movie when does not exist
	movieId := "1"
	movieToUpdate := model.Movie{MovieId: movieId, MovieName: "name1"}
	body, jsonErr := json.Marshal(movieToUpdate)

	req, reqErr := http.NewRequest(http.MethodPut, apiHost+"/movies/1", bytes.NewBuffer(body))
	if jsonErr != nil {
		log.Fatalf("got error when trying to create API request. Error: %s", reqErr)
	}
	req.Header.Set("content-type", "application/json")

	response, err := http.DefaultClient.Do(req)

	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusNotFound, response.StatusCode, "Expected status to be not found")

	responseData, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	}
	responseMessage := model.ResponseMessage{}
	expectedMessage := model.ResponseMessage{Message: "No movie with provided id exists"}
	jsonErr = json.Unmarshal(responseData, &responseMessage)
	if jsonErr != nil {
		log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for not found")

	// Try to create already existing movie
	//clearDatabase()
	//createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})
	//
	//response, err = http.Post(apiHost+"/movies", "application/json", bytes.NewBuffer(body))
	//s.NoErrorf(err, "Should get no error from request initially")
	//s.EqualValuesf(http.StatusConflict, response.StatusCode, "Expected status to be conflict")
	//
	//responseData, readErr = io.ReadAll(response.Body)
	//if readErr != nil {
	//	log.Fatalf("got error when trying to read API response. Error: %s", readErr)
	//}
	//defer response.Body.Close()
	//responseMessage := model.ResponseMessage{}
	//expectedMessage := model.ResponseMessage{false, "A movie with the provided id already exists"}
	//jsonErr = json.Unmarshal(responseData, &responseMessage)
	//if jsonErr != nil {
	//	log.Fatalf("Got error when parsing response. error: %s", jsonErr)
	//}
	//s.Equal(expectedMessage, responseMessage, "Should return message for conflict")
	//
	clearDatabase()
}

func createMovieInDatabase(movie model.Movie) {
	sqlStatement := `INSERT INTO "movies" (movieId, movieName) VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, movie.MovieId, movie.MovieName)
	if err != nil {
		log.Fatalf("Failed creating movie in db for testing: %s", err)
	}
}

func getMovieFromDatabase(movieId string) model.Movie {
	//sqlStatement := `SELECT FROM "movies" WHERE movieId = $1
	movie := model.Movie{}
	err := db.QueryRow("SELECT movieId, movieName FROM movies WHERE movieID = $1;", movieId).
		Scan(&movie.MovieId, &movie.MovieName)

	if err != nil {
		log.Fatalf("Got error when querying db. error: %s", err)
	}

	return movie
}

func clearDatabase() {
	sqlStatement := `DELETE FROM "movies";`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed creating movie in db for testing: %s", err)
	}
}
