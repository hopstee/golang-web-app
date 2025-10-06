package repository

import (
	"errors"
)

var ErrNotFound = errors.New("record not found")

type Repository interface {
	Auth() AuthRepository
	User() UserRepository
	Post() PostRepository
	Request() RequestRepository
}
