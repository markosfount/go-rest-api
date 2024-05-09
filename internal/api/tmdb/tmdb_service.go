package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rest_api/internal/api/config"
	"strings"
)

const (
	tmdbApiUrl  = "https://api.themoviedb.org/3/search/movie"
	apiKeyParam = "?api_key=%s"
	titleQuery  = "&query=%s"
	idQuery     = "/%d"
)

type Service struct {
	Client *http.Client
}

func (s *Service) GetMovieByTitle(title string) (*Movie, error) {
	url := strings.Join([]string{tmdbApiUrl, fmt.Sprintf(apiKeyParam, config.API_KEY), fmt.Sprintf(titleQuery, title)}, "")
	response, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	responseBytes, err := io.ReadAll(response.Body)
	moviesResponse := GetMoviesResponse{}
	err = json.Unmarshal(responseBytes, &moviesResponse)
	if err != nil {
		return nil, err
	}
	if moviesResponse.TotalResults == 0 {
		return nil, ErrNoMoviesFound
	}
	movieId := &moviesResponse.Results[0].ID
	ex := fmt.Sprintf(idQuery, movieId)
	fmt.Println(ex)
	ex2 := fmt.Sprintf("%d", movieId)
	fmt.Println(ex2)
	byIDUrl := strings.Join([]string{tmdbApiUrl, fmt.Sprintf(idQuery, movieId), fmt.Sprintf(apiKeyParam, config.API_KEY)}, "")

	response, err = s.Client.Get(byIDUrl)
	if err != nil {
		return nil, err
	}
	responseBytes, err = io.ReadAll(response.Body)
	movie := &Movie{}
	err = json.Unmarshal(responseBytes, movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *Service) GetMovieByID(id int) (*Movie, error) {
	url := strings.Join([]string{tmdbApiUrl, fmt.Sprintf(idQuery, id), fmt.Sprintf(apiKeyParam, config.API_KEY)}, "")
	response, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	responseBytes, err := io.ReadAll(response.Body)
	movie := &Movie{}
	err = json.Unmarshal(responseBytes, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}
