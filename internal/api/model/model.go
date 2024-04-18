package model

type Movie struct {
	MovieId   string `json:"id"`
	MovieName string `json:"title"`
}

type User struct {
	Username string
	Password string
}

type ResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
