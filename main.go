package main // Update the AllBooks function so it accepts the connection pool as a

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rest_api/handler"
)

const (
	DB_USER     = "user"
	DB_PASSWORD = "password"
	DB_NAME     = "postgres"
)

// DB set up
func setupDB() *sql.DB {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		log.Fatalf("test init failed: %s", err)
	}

	return db
}

func main() {
	// Initialise the connection pool.
	db := setupDB()

	// Create an instance of Env containing the connection pool.
	env := &handler.Env{Db: db}

	http.HandleFunc("/movies", env.GetMovies)
	http.HandleFunc("/movies/create", env.AddMovie)
	http.HandleFunc("/ping", env.TestHandler)

	// listen port
	err := http.ListenAndServe(":3000", nil)
	// print any server-based error messages
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
