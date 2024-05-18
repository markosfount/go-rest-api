package service

import (
	"errors"
	"log/slog"
	"rest_api/internal/api/model"
	"rest_api/internal/data"
)

type MovieService struct {
	movieRepository data.Repository[*model.Movie]
}

func NewMovieService(repository data.Repository[*model.Movie]) MovieService {
	return MovieService{movieRepository: repository}
}

func (s *MovieService) GetAll() ([]*model.Movie, error) {
	movies, err := s.movieRepository.GetAll()
	if err != nil {
		slog.Error("Error when getting movies from db: %s\n", err)
		return []*model.Movie{}, err
	}
	return movies, nil
}

func (s *MovieService) Get(movieId int) (*model.Movie, error) {
	movie, err := s.movieRepository.Get(movieId)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil, model.NotFoundError{}
		}
		slog.Error("Error when getting movie with id %s from db: %s\n", err, movieId)
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) Create(movie *model.Movie) (*model.Movie, error) {
	createdMovie, err := s.movieRepository.Create(movie)

	if err != nil {
		if errors.Is(err, data.ErrRecordExists) {
			return nil, model.ConflictError{}
		}
		slog.Error("Unable to create movie in the database: error: %s\n", err)
		return nil, err
	}
	return createdMovie, nil
}

func (s *MovieService) Update(movie *model.Movie) (*model.Movie, error) {
	updatedMovie, err := s.movieRepository.Update(movie)

	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil, model.NotFoundError{}
		}
		slog.Error("Error when updating movie with id %s in the db: %s\n", err, movie.MovieId)
		return nil, err
	}
	return updatedMovie, nil
}

func (s *MovieService) Delete(movieId int) error {
	err := s.movieRepository.Delete(movieId)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil
		}
		slog.Error("Error when deleting movie in db: %s\n", err)
		return err
	}
	return nil
}
