package repository

import "time"

type Request struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
	ContactType string    `json:"contact_type,omitempty"`
	Message     string    `json:"message"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RequestRepository interface {
	Get() ([]Request, error)
	GetByID(id int64) (Request, error)
	Create(request Request) (int64, error)
	Update(request Request) error
	Delete(id int64) error
}
