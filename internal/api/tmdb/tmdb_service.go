package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const TMDB_API_URL = "https://api.themoviedb.org/3/search/movie?query=%s&api_key=%s"

type Service struct {
	Client *http.Client
}

func (s *Service) GetMovieInfo(title string) (*Movie, error) {
	response, err := s.Client.Get(fmt.Sprintf(TMDB_API_URL, title, os.Getenv("API_KEY")))
	if err != nil {
		return nil, err
	}
	//todo hanlde 404 etc
	responseBytes, err := io.ReadAll(response.Body)

	movie := &Movie{}
	err = json.Unmarshal(responseBytes, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}
