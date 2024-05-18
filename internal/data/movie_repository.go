package data

import (
	"database/sql"
	"errors"
	"fmt"
	"rest_api/internal/api/model"

	"github.com/lib/pq"
)

type MovieRepository struct {
	DB *sql.DB
}

// TODO handle not found in all
func (r *MovieRepository) GetAll() ([]*model.Movie, error) {
	fmt.Println("Getting movies...")

	rows, err := r.DB.Query("SELECT movieId, movieName, overview FROM movies")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*model.Movie

	for rows.Next() {
		var movieID int
		var movieName string
		var overview sql.NullString

		err := rows.Scan(&movieID, &movieName, &overview)
		if err != nil {
			return nil, err
		}
		movie := &model.Movie{
			MovieId:   movieID,
			MovieName: movieName,
		}
		if overview.Valid {
			movie.Overview = overview.String
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) Get(movieId int) (*model.Movie, error) {
	fmt.Printf("Getting movie with movieId %d\n", movieId)

	var overview sql.NullString
	movie := model.Movie{}
	err := r.DB.QueryRow("SELECT movieId, movieName, overview FROM movies WHERE movieID = $1;", movieId).
		Scan(&movie.MovieId, &movie.MovieName, &overview)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	if overview.Valid {
		movie.Overview = overview.String
	}

	return &movie, nil
}

func (r *MovieRepository) Create(movie *model.Movie) (*model.Movie, error) {
	fmt.Printf("Inserting new movie with ID: %d  and name: %s\n", movie.MovieId, movie.MovieName)

	var lastInsertID int
	err := r.DB.QueryRow(
		"INSERT INTO movies(movieID, movieName, overview) VALUES($1, $2, $3) returning id;", movie.MovieId, movie.MovieName, movie.Overview).Scan(&lastInsertID)
	if err != nil {
		pqErr := err.(*pq.Error)
		switch pqErr.Code {
		case "23505":
			return nil, ErrRecordExists
		default:
			return nil, err
		}
	}

	return movie, nil
}

func (r *MovieRepository) Update(movie *model.Movie) (*model.Movie, error) {
	fmt.Printf("Updating movie with ID: %d\n", movie.MovieId)

	res, err := r.DB.Exec(
		"UPDATE movies SET movieId = $1, movieName = $2, overview = $3 WHERE movieId = $1;", movie.MovieId, movie.MovieName, movie.Overview)

	if err != nil {
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrRecordNotFound
	}
	return movie, nil
}

func (r *MovieRepository) Delete(movieId int) error {
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
		return ErrRecordNotFound
	}
	return nil
}
