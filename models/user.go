package models

import (
	"errors"
	"events_booking/db"
	"events_booking/utilis"
	"strings"
)

var ErrUserExists = errors.New("user already exists")

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (user *User) Save() error {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utilis.HashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	result, err := stmt.Exec(user.Email, user.Password)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return ErrUserExists
		}
		return err
	}

	userId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	user.Id = userId
	return nil
}
func (user *User) ValidateCredentials() error {
	query := `SELECT password, id FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, user.Email)

	var storedHashedPassword string

	err := row.Scan(&storedHashedPassword, &user.Id)

	if err != nil {
		return err
	}

	passwordIsValid := utilis.CheckPasswordHash(user.Password, storedHashedPassword)

	if !passwordIsValid {
		return errors.New("invalid credentials unauthorized")
	}

	return nil
}
