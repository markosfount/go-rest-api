package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"rest_api/model"
	"rest_api/utils"
)

type Env struct {
	Db *sql.DB
}

func (env Env) TestHandler(res http.ResponseWriter, req *http.Request) {
	HandlerMessage := []byte(`{
	"success": true,
	"message": "The server is running properly"
}`)

	utils.ReturnJsonResponse(res, http.StatusOK, HandlerMessage)
}

// TODO log errors
func (env Env) GetMovies(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		// Add the response return message
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Check your HTTP method: Invalid HTTP method executed",
}`)

		utils.ReturnJsonResponse(res, http.StatusMethodNotAllowed, HandlerMessage)
		return
	}

	movies, err := model.GetMovies(env.Db)
	if err != nil {
		log.Print(err)
		http.Error(res, http.StatusText(500), 500)
	}

	movieJSON, err := json.Marshal(&movies)
	if err != nil {
		// Add the response return message
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Error parsing the movie data",
}`)

		utils.ReturnJsonResponse(res, http.StatusInternalServerError, HandlerMessage)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

func (env Env) GetMovie(res http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		// Add the response return message
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Check your HTTP method: Invalid HTTP method executed",
}`)

		utils.ReturnJsonResponse(res, http.StatusMethodNotAllowed, HandlerMessage)
		return
	}

	if _, ok := req.URL.Query()["id"]; !ok {
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Μovie id not provided",
}`)
		utils.ReturnJsonResponse(res, http.StatusBadRequest, HandlerMessage)
		return
	}

	id := req.URL.Query()["id"][0]

	movie, err := model.GetMovie(env.Db, id)
	if err != nil {
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Requested movie not found",
}`)

		utils.ReturnJsonResponse(res, http.StatusNotFound, HandlerMessage)
		return
	}

	movieJSON, err := json.Marshal(&movie)
	if err != nil {
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Error parsing the movie data",
}`)

		utils.ReturnJsonResponse(res, http.StatusInternalServerError, HandlerMessage)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

// TODO conflict same movieid
func (env Env) AddMovie(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		// Add the response return message
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Check your HTTP method: Invalid HTTP method executed",
}`)

		utils.ReturnJsonResponse(res, http.StatusMethodNotAllowed, HandlerMessage)
		return
	}

	var movie model.Movie

	payload := req.Body

	defer req.Body.Close()
	err := json.NewDecoder(payload).Decode(&movie)
	if err != nil {
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Error parsing the movie data",
}`)

		utils.ReturnJsonResponse(res, http.StatusInternalServerError, HandlerMessage)
		return
	}

	if movie.MovieId == "" || movie.MovieName == "" {
		HandlerMessage := []byte(`{
	"success": false,
	"message": ""You are missing movieID or movieName parameter",
}`)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, HandlerMessage)
		return
	}
	createdMovie, err := model.CreateMovie(&movie, env.Db)

	if err != nil {
		var cerr *model.ConflictError
		if errors.As(err, &cerr) {
			HandlerMessage := []byte(`{
	"success": false,
	"message": "A movie with the provided id already exists",
}`)
			utils.ReturnJsonResponse(res, http.StatusConflict, HandlerMessage)
			return
		}
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Unexpected error when creating response",
}`)
		fmt.Printf("Unable to create movie in the database: error: %v\n", err)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, HandlerMessage)
		return
	}

	movieJSON, err := json.Marshal(createdMovie)
	if err != nil {
		HandlerMessage := []byte(`{
	"success": false,
	"message": "Unexpected error when creating response",
}`)
		fmt.Printf("Unable to parse movie dao to json: error: %v\n", err)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, HandlerMessage)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusCreated, movieJSON)
}
