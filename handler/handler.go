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
	// Add the response return message
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

	// We can now access the connection pool directly in our handlers.
	movies, err := model.GetMovies(env.Db)
	if err != nil {
		log.Print(err)
		http.Error(res, http.StatusText(500), 500)
	}

	for _, movie := range movies {
		fmt.Fprintf(res, "%s, %s", movie.MovieId, movie.MovieName)
	}

	// parse the movie data into json format
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
