package repository

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password,omitempty"`
	Email        string    `json:"email,omitempty"`
	TokenVersion int       `json:"token_version"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthData struct {
	ID           int64  `json:"id"`
	Password     string `json:"password"`
	TokenVersion int    `json:"token_version"`
}

type UserRepository interface {
	GetByID(id int64) (User, error)
	GetByUsername(username string) (User, error)
	GetByEmail(email string) (User, error)
	GetAuthData(username string) (AuthData, error)
	GetTokenVersion(id int64) (int, error)
	Create(user User) (int64, error)
	Update(user User) error
	Delete(id int64) error
	IncrementTokenVersion(id int64) error
}
