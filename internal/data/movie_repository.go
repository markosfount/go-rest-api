package data

import (
	"database/sql"
	"fmt"
	"rest_api/internal/api/model"

	"github.com/lib/pq"
)

type MovieRepository struct {
	DB *sql.DB
}

// TODO handle not found in all
func (r *MovieRepository) GetMovies() ([]model.Movie, error) {
	fmt.Println("Getting movies...")

	rows, err := r.DB.Query("SELECT * FROM movies")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []model.Movie

	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err := rows.Scan(&id, &movieID, &movieName)
		if err != nil {
			return nil, err
		}

		movies = append(movies, model.Movie{MovieId: movieID, MovieName: movieName})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) GetMovie(movieId string) (*model.Movie, error) {
	fmt.Printf("Getting movie with movieId %s\n", movieId)

	movie := model.Movie{}
	err := r.DB.QueryRow("SELECT movieId, movieName FROM movies WHERE movieID = $1;", movieId).
		Scan(&movie.MovieId, &movie.MovieName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &model.NotFoundError{}
		}
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) CreateMovie(movie *model.Movie) (*model.Movie, error) {
	fmt.Println("Inserting new movie with ID: " + movie.MovieId + " and name: " + movie.MovieName)

	var lastInsertID int
	err := r.DB.QueryRow(
		"INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movie.MovieId, movie.MovieName).Scan(&lastInsertID)
	if err != nil {
		pqErr := err.(*pq.Error)
		switch pqErr.Code {
		case "23505":
			return nil, &model.ConflictError{}
		default:
			return nil, err
		}
	}

	return movie, nil
}

func (r *MovieRepository) UpdateMovie(movie *model.Movie) (*model.Movie, error) {
	fmt.Println("Updating movie with ID: " + movie.MovieId)

	res, err := r.DB.Exec(
		"UPDATE movies SET movieId = $1, movieName = $2 WHERE movieId = $1;", movie.MovieId, movie.MovieName)

	if err != nil {
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, &model.NotFoundError{}
	}
	return movie, nil
}

func (r *MovieRepository) DeleteMovie(movieId string) error {
	fmt.Printf("Deleting movie with movieId %s\n", movieId)

	res, err := r.DB.Exec("DELETE FROM movies WHERE movieID = $1;", movieId)

	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return &model.NotFoundError{}
	}
	return nil
}