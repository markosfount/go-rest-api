package model

type Movie struct {
	MovieId   int    `json:"id"`
	MovieName string `json:"title"`
	Overview  string `json:"overview"`
}

type User struct {
	Username string
	Password string
}

type ResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
