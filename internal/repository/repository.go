package repository

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("record not found")

type Repository interface {
	Database() *sql.DB

	Auth() AuthRepository
	User() UserRepository
	Admin() AdminRepository
	Post() PostRepository
	Request() RequestRepository
	Pages() PagesRepository
}
