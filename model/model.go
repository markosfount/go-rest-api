package model

import (
	"database/sql"
	"fmt"
)

type Movie struct {
	MovieId   string `json:"id"`
	MovieName string `json:"title"`
}

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

func CreateMovie(movie *Movie, db *sql.DB) (*Movie, error) {
	fmt.Println("Inserting new movie with ID: " + movie.MovieId + " and name: " + movie.MovieName)

	var lastInsertID int
	err := db.QueryRow(
		"INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movie.MovieId, movie.MovieName).Scan(&lastInsertID)

	if err != nil {
		return nil, err
	}

	return movie, nil
}
