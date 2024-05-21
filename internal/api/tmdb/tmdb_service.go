package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rest_api/internal/api/config"
	"strings"
)

const (
	apiUrl          = "https://api.themoviedb.org/3/"
	searchEndpoint  = "search/movie"
	detailsEndpoint = "/movie"
	apiKeyParam     = "?api_key=%s"
	titleQuery      = "&query=%s"
	idQuery         = "/%d"
)

type Service struct {
	client *http.Client
}

func NewService() Service {
	return Service{client: http.DefaultClient}
}

func (s *Service) GetMovieByTitle(title string) (*Movie, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, createSearchUrl(title), nil)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Do(request)
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
	//movieId := moviesResponse.Results[0].ID
	movie := moviesResponse.Results[0]

	//response, err = s.client.Get(createDetailsUrl(movieId))
	//if err != nil {
	//	return nil, err
	//}
	//responseBytes, err = io.ReadAll(response.Body)
	//movie := &Movie{}
	err = json.Unmarshal(responseBytes, movie)
	//if err != nil {
	//	return nil, err
	//}

	//return movie, nil
	return &movie, nil
}

func (s *Service) GetMovieByID(id int) (*Movie, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, createDetailsUrl(id), nil)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Do(request)
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

func createSearchUrl(title string) string {
	return strings.Join([]string{apiUrl, searchEndpoint, fmt.Sprintf(apiKeyParam, config.ApiKey), fmt.Sprintf(titleQuery, title)}, "")
}

func createDetailsUrl(id int) string {
	return strings.Join([]string{apiUrl, detailsEndpoint, fmt.Sprintf(idQuery, id), fmt.Sprintf(apiKeyParam, config.ApiKey)}, "")
}
