package tmdb

import "errors"

var ErrNoMoviesFound = errors.New("no movies found")

type Movie struct {
	ID            int    `json:"id"`
	OriginalTitle string `json:"original_title"`
	Overview      string `json:"overview"`
	Runtime       int32  `json:"runtime"`
}

type GetMoviesResponse struct {
	Results      []Movie `json:"results"`
	TotalResults int     `json:"total_results""`
}
