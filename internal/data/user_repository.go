package data

import (
	"database/sql"
	"rest_api/internal/api/model"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetUser(username string) (*model.User, error) {
	user := model.User{}
	err := r.DB.QueryRow("SELECT username, password FROM users WHERE username = $1;", username).
		Scan(user.Username, user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &model.NotFoundError{}
		}
		return nil, err
	}

	return &user, nil
}