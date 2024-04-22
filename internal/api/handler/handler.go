package handler

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"rest_api/internal/api/model"
	"rest_api/internal/api/service"
	"rest_api/internal/api/utils"
	"rest_api/internal/data"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Handler struct {
	UserRepository data.UserRepository
	MovieService   service.MovieService
}

func (h *Handler) PingHandler(res http.ResponseWriter, _ *http.Request) {
	responseBytes := createResponse(true, "The server is running properly")
	utils.ReturnJsonResponse(res, http.StatusOK, responseBytes)
}

func (h *Handler) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if ok {
			//usernameHash := sha256.Sum256([]byte(username))
			user, err := h.UserRepository.GetUser(username)
			if err != nil {
				var nferr *model.NotFoundError
				if errors.As(err, &nferr) {
					utils.ReturnUnauthorizedResponse(res)
					return
				}
				log.Printf("Error when getting user from db: %s\n", err)
				responseBytes := createResponse(false, "Error when accessing user list. Please try again")
				utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
				return
			}
			passwordHash := sha256.Sum256([]byte(password))
			//expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(user.Password))

			//usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if passwordMatch {
				next.ServeHTTP(res, req)
				return
			}
		}
		utils.ReturnUnauthorizedResponse(res)
	}
}

func (h *Handler) GetMovies(res http.ResponseWriter, _ *http.Request) {
	movies, err := h.MovieService.GetAll()
	if err != nil {
		returnErrorResponse("Error when retrieving data", http.StatusInternalServerError, res)
	}
	movieJSON, err := json.Marshal(&movies)
	if err != nil {
		slog.Error("Error when marshalling the response data: %s\n", err)
		returnErrorResponse("Error creating response", http.StatusInternalServerError, res)
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

func (h *Handler) GetMovie(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	movieId := vars["movieId"]

	movie, err := h.MovieService.MovieRepository.Get(movieId)
	if err != nil {
		var nfErr *model.NotFoundError
		if errors.As(err, &nfErr) {
			returnErrorResponse("No movie with provided id exists", http.StatusNotFound, res)
			return
		}
		responseBytes := createResponse(false, "Error when retrieving requested movie")
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	movieJSON, err := json.Marshal(&movie)
	if err != nil {
		slog.Error("Error when marshalling the response data: %s\n", err)
		returnErrorResponse("Error creating response", http.StatusInternalServerError, res)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

func (h *Handler) AddMovie(res http.ResponseWriter, req *http.Request) {
	var movie *model.Movie

	payload := req.Body

	defer req.Body.Close()
	err := json.NewDecoder(payload).Decode(&movie)
	if err != nil {
		slog.Error("Error when unmarshalling the request body: %s\n", err)
		returnErrorResponse("Could not parse request body", http.StatusBadRequest, res)
		return
	}

	if movie.MovieId == "" || movie.MovieName == "" {
		returnErrorResponse("Missing movieID or movieName parameter", http.StatusBadRequest, res)
		return
	}
	createdMovie, err := h.MovieService.Create(movie)

	if err != nil {
		var cErr *model.ConflictError
		if errors.As(err, &cErr) {
			returnErrorResponse("A movie with the provided id already exists", http.StatusConflict, res)
			return
		}
		returnErrorResponse("Unexpected error when creating data", http.StatusConflict, res)
		return
	}

	movieJSON, err := json.Marshal(createdMovie)
	if err != nil {
		slog.Error("Error when marshalling the response data: %s\n", err)
		returnErrorResponse("Error creating response", http.StatusInternalServerError, res)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusCreated, movieJSON)
}

func (h *Handler) UpdateMovie(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	movieId := vars["movieId"]

	var movie *model.Movie

	payload := req.Body

	defer req.Body.Close()
	err := json.NewDecoder(payload).Decode(&movie)
	if err != nil {
		slog.Error("Error when unmarshalling the request body: %s\n", err)
		returnErrorResponse("Could not parse request body", http.StatusInternalServerError, res)
		return
	}

	if movieId != movie.MovieId {
		returnErrorResponse("Mismatch between movieId in query parameter and request body", http.StatusBadRequest, res)
		return
	}
	updatedMovie, err := h.MovieService.Update(movie)

	if err != nil {
		var nfErr *model.NotFoundError
		if errors.As(err, &nfErr) {
			responseBytes := createResponse(false, "No movie with provided id exists")
			utils.ReturnJsonResponse(res, http.StatusNotFound, responseBytes)
			return
		}
		responseBytes := createResponse(false, "Unexpected error when updating movie.")
		utils.ReturnJsonResponse(res, http.StatusInternalServerError, responseBytes)
		return
	}

	movieJSON, err := json.Marshal(updatedMovie)
	if err != nil {
		slog.Error("Error when marshalling the response data: %s\n", err)
		returnErrorResponse("Error creating response", http.StatusInternalServerError, res)
		return
	}

	utils.ReturnJsonResponse(res, http.StatusOK, movieJSON)
}

func (h *Handler) DeleteMovie(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	movieId := vars["movieId"]

	err := h.MovieService.Delete(movieId)
	if err != nil {
		returnErrorResponse("Error when deleting requested movie", http.StatusInternalServerError, res)
		return
	}

	utils.ReturnEmptyResponse(res, http.StatusNoContent)
}

func returnErrorResponse(message string, status int, res http.ResponseWriter) {
	responseBytes := createResponse(false, message)
	utils.ReturnJsonResponse(res, status, responseBytes)
}

func createResponse(success bool, message string) []byte {
	response := model.ResponseMessage{Success: success, Message: message}
	responseBytes, _ := json.Marshal(response)
	return responseBytes
}
