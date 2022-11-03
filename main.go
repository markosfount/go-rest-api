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

// Create a custom Env struct which holds a connection pool.
type Env struct {
	db *sql.DB
}

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatalf("test init failed: %s", err)
	}

	return db
}

func main() {
	// Initialise the connection pool.
	db := setupDB()

	// Create an instance of Env containing the connection pool.
	env := &Env{db: db}

	http.HandleFunc("/", handler.TestHandler)

	// listen port
	err := http.ListenAndServe(":3000", nil)
	// print any server-based error messages
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Use env.booksIndex as the handler function for the /books route.
	http.HandleFunc("/movies", env.booksIndex)
	http.ListenAndServe(":3000", nil)
}

// Define booksIndex as a method on Env.
func (env *Env) booksIndex(w http.ResponseWriter, r *http.Request) {
	// We can now access the connection pool directly in our handlers.
	bks, err := models.AllBooks(env.db)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.Isbn, bk.Title, bk.Author, bk.Price)
	}
}
