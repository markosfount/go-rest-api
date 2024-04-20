package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rest_api/internal/api/application"
	"rest_api/internal/api/model"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

const (
	databaseName = "movies"
	databasePort = 5432
)

var databaseHost = application.GetEnv("DB_HOST", "localhost")
var apiHost = application.GetEnv("API_HOST", "http://localhost:3000")

var databaseUser = application.GetEnv("DB_USER", "root")
var databasePassword = application.GetEnv("DB_PASS", "password")

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
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be ok")

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	var movies []model.Movie
	err = json.Unmarshal(responseData, &movies)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.ElementsMatchf([]model.Movie{}, movies, "Should return empty list")

	// GET when movies exist in db
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})
	createMovieInDatabase(model.Movie{MovieId: "2", MovieName: "name2"})

	response, err = http.Get(apiHost + "/movies")
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be ok")

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	movies = []model.Movie{}
	expectedMovies := []model.Movie{{MovieId: "1", MovieName: "name1"}, {MovieId: "2", MovieName: "name2"}}
	err = json.Unmarshal(responseData, &movies)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.ElementsMatchf(expectedMovies, movies, "Should return two movies")

	clearDatabase()
}

// func (s *AppSuite) TestAuthentication() {
// 	// GET when no authorization provided

// 	response, err := http.Get(apiHost + "/movies")
// 	s.NoErrorf(err, "Should get no error from request initally")
// 	s.EqualValuesf(http.StatusUnauthorized, response.StatusCode, "Expected status to be: Unauthorized")

// 	responseData, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatalf("got error when trying to read API response. Error: %s", err)
// 	}
// 	defer response.Body.Close()
// 	err = json.Unmarshal(responseData, &movies)
// 	if err != nil {
// 		log.Fatalf("Got error when parsing response. error: %s", err)
// 	}
// 	s.ElementsMatchf([]model.Movie{}, movies, "Should return empty list")

// 	// GET when movies exist in db
// 	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})
// 	createMovieInDatabase(model.Movie{MovieId: "2", MovieName: "name2"})

// 	response, err = http.Get(apiHost + "/movies")
// 	s.NoErrorf(err, "Should get no error from request initially")
// 	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be ok")

// 	responseData, err = io.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatalf("got error when trying to read API response. Error: %s", err)
// 	}
// 	defer response.Body.Close()
// 	movies = []model.Movie{}
// 	expectedMovies := []model.Movie{{MovieId: "1", MovieName: "name1"}, {MovieId: "2", MovieName: "name2"}}
// 	err = json.Unmarshal(responseData, &movies)
// 	if err != nil {
// 		log.Fatalf("Got error when parsing response. error: %s", err)
// 	}
// 	s.ElementsMatchf(expectedMovies, movies, "Should return two movies")

// 	clearDatabase()
// }

func (s *AppSuite) TestGetMovie() {
	// GET when movie does not exist
	response, err := http.Get(apiHost + "/movies/1")
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusNotFound, response.StatusCode, "Expected status to be not found")

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	responseMessage := model.ResponseMessage{}
	expectedMessage := model.ResponseMessage{Message: "No movie with provided id exists"}
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for not found")

	// GET when movies exist in db
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})

	response, err = http.Get(apiHost + "/movies/1")
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be ok")

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	movie := model.Movie{}
	expectedMovie := model.Movie{MovieId: "1", MovieName: "name1"}
	err = json.Unmarshal(responseData, &movie)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMovie, movie, "Should return created movie")

	clearDatabase()
}

func (s *AppSuite) TestCreateMovie() {
	// Create movie
	movieId := "1"
	movieToCreate := model.Movie{MovieId: movieId, MovieName: "name1"}
	body, err := json.Marshal(movieToCreate)

	response, err := http.Post(apiHost+"/movies", "application/json", bytes.NewBuffer(body))
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusCreated, response.StatusCode, "Expected status to be: Created")

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	movie := model.Movie{}
	expectedMovie := model.Movie{MovieId: "1", MovieName: "name1"}
	err = json.Unmarshal(responseData, &movie)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMovie, movie, "Should return created movie")
	// check that movie was created in db
	savedMovie := getMovieFromDatabase(movieId)
	s.Equal(expectedMovie, savedMovie, "Should return created movie")

	// Create movie with malformed request
	body = []byte("{asf}")

	response, err = http.Post(apiHost+"/movies", "application/json", bytes.NewBuffer(body))
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusBadRequest, response.StatusCode, "Expected status to be: Bad Request")

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	responseMessage := model.ResponseMessage{}
	expectedMessage := model.ResponseMessage{Message: "Could not parse request body"}
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for: Bad Request")
	clearDatabase()

	// Create movie with missing id in body
	movieToCreate = model.Movie{MovieName: "name1"}
	body, err = json.Marshal(movieToCreate)

	response, err = http.Post(apiHost+"/movies", "application/json", bytes.NewBuffer(body))
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusBadRequest, response.StatusCode, "Expected status to be: Bad Request")

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	responseMessage = model.ResponseMessage{}
	expectedMessage = model.ResponseMessage{Message: "You are missing movieID or movieName parameter"}
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for: Bad Request")
	clearDatabase()

	// Try to create already existing movie
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})
	movieToCreate = model.Movie{MovieId: movieId, MovieName: "name2"}
	body, err = json.Marshal(movieToCreate)

	response, err = http.Post(apiHost+"/movies", "application/json", bytes.NewBuffer(body))
	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusConflict, response.StatusCode, "Expected status to be: Conflict")

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	responseMessage = model.ResponseMessage{}
	expectedMessage = model.ResponseMessage{Message: "A movie with the provided id already exists"}
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for: Conflict")

	clearDatabase()
}

func (s *AppSuite) TestUpdateMovie() {
	// Update movie when does not exist
	movieId := "1"
	movieToUpdate := model.Movie{MovieId: movieId, MovieName: "name1"}
	body, err := json.Marshal(movieToUpdate)

	req, err := http.NewRequest(http.MethodPut, apiHost+"/movies/1", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("got error when trying to create API request. Error: %s", err)
	}
	req.Header.Set("content-type", "application/json")

	response, err := http.DefaultClient.Do(req)

	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusNotFound, response.StatusCode, "Expected status to be: Not found")

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	responseMessage := model.ResponseMessage{}
	expectedMessage := model.ResponseMessage{Message: "No movie with provided id exists"}
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for: Not found")

	// Mismatch between path param and body
	movieToUpdate = model.Movie{MovieId: movieId, MovieName: "name1"}
	body, err = json.Marshal(movieToUpdate)

	req, err = http.NewRequest(http.MethodPut, apiHost+"/movies/2", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("got error when trying to create API request. Error: %s", err)
	}
	req.Header.Set("content-type", "application/json")

	response, err = http.DefaultClient.Do(req)

	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusBadRequest, response.StatusCode, "Expected status to be: Bad request")

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	responseMessage = model.ResponseMessage{}
	expectedMessage = model.ResponseMessage{Message: "Mismatch between movieId in query parameter and request body"}
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for mismatch")

	clearDatabase()

	//Update existing movie
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})
	movieToUpdate = model.Movie{MovieId: "1", MovieName: "name2"}
	body, err = json.Marshal(movieToUpdate)

	req, err = http.NewRequest(http.MethodPut, apiHost+"/movies/1", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("got error when trying to create API request. Error: %s", err)
	}
	req.Header.Set("content-type", "application/json")

	response, err = http.DefaultClient.Do(req)

	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusOK, response.StatusCode, "Expected status to be: Ok")

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}

	movie := model.Movie{}
	expectedMovie := model.Movie{MovieId: "1", MovieName: "name2"}
	err = json.Unmarshal(responseData, &movie)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMovie, movie, "Should return updated movie")

	// check that movie was updated in db
	updatedMovie := getMovieFromDatabase(movieId)
	s.Equal(expectedMovie, updatedMovie, "Should return updated movie")

	clearDatabase()
}

func (s *AppSuite) TestDeleteMovie() {
	// DELETE when movie does not exist
	req, err := http.NewRequest(http.MethodDelete, apiHost+"/movies/1", nil)
	if err != nil {
		log.Fatalf("got error when trying to create API request. Error: %s", err)
	}
	response, err := http.DefaultClient.Do(req)

	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusNotFound, response.StatusCode, "Expected status to be not found")

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("got error when trying to read API response. Error: %s", err)
	}
	defer response.Body.Close()
	responseMessage := model.ResponseMessage{}
	expectedMessage := model.ResponseMessage{Message: "No movie with provided id exists"}
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		log.Fatalf("Got error when parsing response. error: %s", err)
	}
	s.Equal(expectedMessage, responseMessage, "Should return message for not found")

	// DELETE when movies exist in db
	createMovieInDatabase(model.Movie{MovieId: "1", MovieName: "name1"})
	req, err = http.NewRequest(http.MethodDelete, apiHost+"/movies/1", nil)
	if err != nil {
		log.Fatalf("got error when trying to create API request. Error: %s", err)
	}

	response, err = http.DefaultClient.Do(req)

	s.NoErrorf(err, "Should get no error from request initially")
	s.EqualValuesf(http.StatusNoContent, response.StatusCode, "Expected status to be: No content")

	clearDatabase()
}

func createMovieInDatabase(movie model.Movie) {
	sqlStatement := `INSERT INTO "movies" (movieId, movieName) VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, movie.MovieId, movie.MovieName)
	if err != nil {
		log.Fatalf("Failed creating movie in db for testing: %s", err)
	}
}

func createUserInDatabase(username, password string) {
	sqlStatement := `INSERT INTO "users" (username, password) VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, username, password)
	if err != nil {
		log.Fatalf("Failed creating user in db for testing: %s", err)
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
	sqlStatement := `DELETE FROM "movies"; DELETE FROM "users";`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed clearing database: %s", err)
	}
}
