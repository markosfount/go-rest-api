package model

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type Movie struct {
	MovieId   string `json:"id"`
	MovieName string `json:"title"`
}

type ResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ConflictError struct {
}

func (c *ConflictError) Error() string {
	return "Conflict when trying to add movie."
}

type NotFoundError struct {
}

func (c *NotFoundError) Error() string {
	return "Movie not found."
}

// TODO handle not found in all
func GetMovies(db *sql.DB) ([]Movie, error) {
	fmt.Println("Getting movies...")

	rows, err := db.Query("SELECT * FROM movies")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err := rows.Scan(&id, &movieID, &movieName)
		if err != nil {
			return nil, err
		}

		movies = append(movies, Movie{MovieId: movieID, MovieName: movieName})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func GetMovie(db *sql.DB, movieId string) (*Movie, error) {
	fmt.Printf("Getting movie with movieId %s\n", movieId)

	movie := Movie{}
	err := db.QueryRow("SELECT movieId, movieName FROM movies WHERE movieID = $1;", movieId).
		Scan(&movie.MovieId, &movie.MovieName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{}
		}
		return nil, err
	}

	return &movie, nil
}

func CreateMovie(movie *Movie, db *sql.DB) (*Movie, error) {
	fmt.Println("Inserting new movie with ID: " + movie.MovieId + " and name: " + movie.MovieName)

	var lastInsertID int
	err := db.QueryRow(
		"INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movie.MovieId, movie.MovieName).Scan(&lastInsertID)

	if err != nil {
		pqErr := err.(*pq.Error)
		switch pqErr.Code {
		case "23505":
			return nil, &ConflictError{}
		default:
			return nil, err
		}
	}

	return movie, nil
}

func UpdateMovie(movie *Movie, db *sql.DB) (*Movie, error) {
	fmt.Println("Updating movie with ID: " + movie.MovieId)

	res, err := db.Exec(
		"UPDATE movies SET movieId = $1, movieName = $2 WHERE movieId = $1;", movie.MovieId, movie.MovieName)

	if err != nil {
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, &NotFoundError{}
	}
	return movie, nil
}
