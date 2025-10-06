package repository

import "time"

type Post struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	HeroImgURL string    `json:"hero_img_url"`
	Content    string    `json:"content"`
	Likes      int64     `json:"likes"`
	IsPublic   bool      `json:"is_public"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostRepository interface {
	GetAllPublic() ([]Post, error)
	GetPublicByID(id int64) (Post, error)
	GetAll() ([]Post, error)
	GetByID(id int64) (Post, error)
	Create(post Post) (int64, error)
	Update(post Post) error
	Delete(id int64) error
}
