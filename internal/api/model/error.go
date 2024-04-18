package model

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