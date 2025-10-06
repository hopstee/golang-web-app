package repository

import "time"

type Admin struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminAuthData struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

type AdminRepository interface {
	GetByID(id int64) (Admin, error)
	GetByUsername(username string) (Admin, error)
	GetByEmail(email string) (Admin, error)
	GetAuthData(username string) (AuthData, error)
	Create(user Admin) (int64, error)
	Update(user Admin) error
	Delete(id int64) error
}
