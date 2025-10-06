package sqlite

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
		SELECT id, title, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
		WHERE is_public = true
	`)
}

func (r *postRepo) GetAll() ([]repository.Post, error) {
	return r.getAllPosts(`
		SELECT id, title, hero_img_url, content, likes, is_public, updated_at, created_at
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

		err := rows.Scan(&p.ID, &p.Title, &p.HeroImgURL, &p.Content, &p.Likes, &p.IsPublic, &p.UpdatedAt, &p.CreatedAt)
		if err != nil {
			r.Logger.Error("failed to scan post", slog.String("post", p.Title), slog.Any("err", err))
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *postRepo) GetByID(id int64) (repository.Post, error) {
	return r.getPostByID(`
		SELECT id, title, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
		WHERE id = ?
	`, id)
}

func (r *postRepo) GetPublicByID(id int64) (repository.Post, error) {
	return r.getPostByID(`
		SELECT id, title, hero_img_url, content, likes, is_public, updated_at, created_at
		FROM posts
		WHERE id = ? AND is_public = true
	`, id)
}

func (r *postRepo) getPostByID(query string, id int64) (repository.Post, error) {
	rows := r.DB.QueryRow(query, id)

	var post repository.Post

	err := rows.Scan(
		&post.ID,
		&post.Title,
		&post.HeroImgURL,
		&post.Content,
		&post.Likes,
		&post.IsPublic,
		&post.UpdatedAt,
		&post.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("post not found", slog.Int64("id", id))
			return post, repository.ErrNotFound
		}
		r.Logger.Error("failed to query post", slog.Int64("id", id), slog.Any("err", err))
		return post, err
	}

	return post, nil
}

func (r *postRepo) Create(post repository.Post) (int64, error) {
	result, err := r.DB.Exec(`
		INSERT INTO posts(title, hero_img_url, content, likes, is_public, updated_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, post.Title, post.HeroImgURL, post.Content, 0, post.IsPublic, time.Now(), time.Now())
	if err != nil {
		r.Logger.Error("failed to create post", slog.String("post", post.Title), slog.Any("err", err))
		return 0, err
	}

	r.Logger.Debug("post created", slog.String("post", post.Title))
	return result.LastInsertId()
}

func (r *postRepo) Update(post repository.Post) error {
	_, err := r.DB.Exec(`
		UPDATE posts
		SET title = ?, hero_img_url = ?, content = ?, likes = ?, is_public = ?, updated_at = ?
		WHERE id = ?
	`, post.Title, post.HeroImgURL, post.Content, post.Likes, post.IsPublic, time.Now(), post.ID)
	if err != nil {
		r.Logger.Error("failed to update post", slog.String("post", post.Title), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("post updated", slog.Int64("post", post.ID))
	return nil
}

func (r *postRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM posts WHERE id = ?`, id)
	if err != nil {
		r.Logger.Error("failed to delete post", slog.Int64("id", id), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("post deleted", slog.Int64("id", id))
	return nil
}
