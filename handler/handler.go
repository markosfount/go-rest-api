package handler

import (
	"database/sql"
	"encoding/json"
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
	}
	createdMovie, err := model.CreateMovie(&movie, env.Db)

	movieJSON, err := json.Marshal(createdMovie)
	if err != nil {
		HandlerMessage := []byte(`{
   				"success": false,
   				"message": "Unexpected error when creating response",
  			}`)
		fmt.Printf("Unable to parse movie dao to json: error: %v\n", err)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, HandlerMessage)
	}

	utils.ReturnJsonResponse(res, http.StatusCreated, movieJSON)
}
