package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
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
	responseBytes := createResponse(true, "The server is running properly")
	utils.ReturnJsonResponse(res, http.StatusOK, responseBytes)
}

func (env Env) GetMovies(res http.ResponseWriter, req *http.Request) {
	movies, err := model.GetMovies(env.Db)
	if err != nil {
		log.Printf("Error when getting movies from db: %s\n", err)
		responseBytes := createResponse(false, "Error when retrieving data")
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	movieJSON, err := json.Marshal(&movies)
	if err != nil {
		log.Printf("Error when marshalling the response data: %s\n", err)
		responseBytes := createResponse(false, "Error creating response")
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

func (env Env) GetMovie(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	movieId := vars["movieId"]

	// fixme specific error for not found
	movie, err := model.GetMovie(env.Db, movieId)
	if err != nil {
		responseBytes := createResponse(false, "No movie with provided id exists")
		utils.ReturnJsonResponse(res, http.StatusNotFound, responseBytes)
		return
	}

	movieJSON, err := json.Marshal(&movie)
	if err != nil {
		log.Printf("Error when marshalling the response data: %s\n", err)
		responseBytes := createResponse(false, "Error creating response")
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

func (env Env) AddMovie(res http.ResponseWriter, req *http.Request) {
	var movie model.Movie

	payload := req.Body

	defer req.Body.Close()
	err := json.NewDecoder(payload).Decode(&movie)
	if err != nil {
		log.Printf("Error when unmarshalling the request body: %s\n", err)
		responseBytes := createResponse(false, "Could not parse request body")
		utils.ReturnJsonResponse(res, http.StatusBadRequest, responseBytes)
		return
	}

	if movie.MovieId == "" || movie.MovieName == "" {
		responseBytes := createResponse(false, "You are missing movieID or movieName parameter")
		utils.ReturnJsonResponse(res, http.StatusBadRequest, responseBytes)
		return
	}
	createdMovie, err := model.CreateMovie(&movie, env.Db)

	if err != nil {
		var cerr *model.ConflictError
		if errors.As(err, &cerr) {
			responseBytes := createResponse(false, "A movie with the provided id already exists")
			utils.ReturnJsonResponse(res, http.StatusConflict, responseBytes)
			return
		}
		responseBytes := createResponse(false, "Unexpected error when creating data")
		log.Printf("Unable to create movie in the database: error: %s\n", err)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	movieJSON, err := json.Marshal(createdMovie)
	if err != nil {
		responseBytes := createResponse(false, "Error creating response")
		log.Printf("Error when marshalling the response data: %s\n", err)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusCreated, movieJSON)
}

func (env Env) UpdateMovie(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	movieId := vars["movieId"]

	var movie model.Movie

	payload := req.Body

	defer req.Body.Close()
	err := json.NewDecoder(payload).Decode(&movie)
	if err != nil {
		log.Printf("Error when unmarshalling the request body: %s\n", err)
		responseBytes := createResponse(false, "Could not parse request body")
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	if movieId != movie.MovieId {
		responseBytes := createResponse(false, "Mismatch between movieId in query parameter and request body")
		utils.ReturnJsonResponse(res, http.StatusBadRequest, responseBytes)
		return
	}
	updatedMovie, err := model.UpdateMovie(&movie, env.Db)

	if err != nil {
		var nferr *model.NotFoundError
		if errors.As(err, &nferr) {
			responseBytes := createResponse(false, "No movie with provided id exists")
			utils.ReturnJsonResponse(res, http.StatusNotFound, responseBytes)
			return
		}
		responseBytes := createResponse(false, "Unexpected error when updating movie.")
		log.Printf("Unable to update movie in the database: error: %s\n", err)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	movieJSON, err := json.Marshal(updatedMovie)
	if err != nil {
		responseBytes := createResponse(false, "Error creating response")
		log.Printf("Error when marshalling the response data: %s\n", err)
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

func createResponse(success bool, message string) []byte {
	response := model.ResponseMessage{Success: success, Message: message}
	responseBytes, _ := json.Marshal(response)
	return responseBytes
}
