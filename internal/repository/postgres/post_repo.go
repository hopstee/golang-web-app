package postgres

import (
	"database/sql"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"time"
)

type postRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewPostRepo(db *sql.DB, logger *slog.Logger) repository.PostRepository {
	return &postRepo{
		DB:     db,
		Logger: logger,
	}
}

func (r *postRepo) GetAllPublic() ([]repository.Post, error) {
	return r.getAllPosts(`
		SELECT id, title, slug, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
		WHERE is_public = TRUE
	`)
}

func (r *postRepo) GetAll() ([]repository.Post, error) {
	return r.getAllPosts(`
		SELECT id, title, slug, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
	`)
}

func (r *postRepo) getAllPosts(query string, args ...interface{}) ([]repository.Post, error) {
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []repository.Post
	for rows.Next() {
		var p repository.Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.HeroImgURL,
			&p.Content,
			&p.Likes,
			&p.IsPublic,
			&p.UpdatedAt,
			&p.CreatedAt,
		)
		if err != nil {
			r.Logger.Error("failed to scan post", slog.String("post", p.Title), slog.Any("err", err))
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *postRepo) GetByID(id int64) (repository.Post, error) {
	return r.getPost(`
		SELECT id, title, slug, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
		WHERE id = $1
	`, id)
}

func (r *postRepo) GetPublicByID(id int64) (repository.Post, error) {
	return r.getPost(`
		SELECT id, title, slug, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
		WHERE id = $1 AND is_public = TRUE
	`, id)
}

func (r *postRepo) GetBySlug(slug string) (repository.Post, error) {
	return r.getPost(`
		SELECT id, title, slug, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
		WHERE slug = $1
	`, slug)
}

func (r *postRepo) getPost(query string, arg interface{}) (repository.Post, error) {
	row := r.DB.QueryRow(query, arg)
	var post repository.Post

	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.HeroImgURL,
		&post.Content,
		&post.Likes,
		&post.IsPublic,
		&post.UpdatedAt,
		&post.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("post not found", slog.Any("arg", arg))
			return post, repository.ErrNotFound
		}
		r.Logger.Error("failed to query post", slog.Any("arg", arg), slog.Any("err", err))
		return post, err
	}

	return post, nil
}

func (r *postRepo) Create(post repository.Post) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`
		INSERT INTO posts (title, slug, hero_img_url, content, likes, is_public, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, post.Title, post.Slug, post.HeroImgURL, post.Content, 0, post.IsPublic, time.Now(), time.Now(),
	).Scan(&id)
	if err != nil {
		r.Logger.Error("failed to create post", slog.String("title", post.Title), slog.Any("err", err))
		return 0, err
	}

	r.Logger.Debug("post created", slog.Int64("id", id), slog.String("title", post.Title))
	return id, nil
}

func (r *postRepo) Update(post repository.Post) error {
	_, err := r.DB.Exec(`
		UPDATE posts
		SET title = $1, slug = $2, hero_img_url = $3, content = $4, likes = $5, is_public = $6, updated_at = $7
		WHERE id = $8
	`, post.Title, post.Slug, post.HeroImgURL, post.Content, post.Likes, post.IsPublic, time.Now(), post.ID)
	if err != nil {
		r.Logger.Error("failed to update post", slog.String("title", post.Title), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("post updated", slog.Int64("id", post.ID))
	return nil
}

func (r *postRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		r.Logger.Error("failed to delete post", slog.Int64("id", id), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("post deleted", slog.Int64("id", id))
	return nil
}
