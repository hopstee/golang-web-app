package sqlite

import (
	"database/sql"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"time"
)

type requestRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewRequestRepo(db *sql.DB, logger *slog.Logger) repository.RequestRepository {
	return &requestRepo{
		DB:     db,
		Logger: logger,
	}
}

func (r *requestRepo) Get() ([]repository.Request, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, phone, email, contact_type, message, created_at, updated_at
		FROM requests
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []repository.Request
	for rows.Next() {
		var req repository.Request

		err := rows.Scan(
			&req.ID,
			&req.Name,
			&req.Phone,
			&req.Email,
			&req.ContactType,
			&req.Message,
			&req.CreatedAt,
			&req.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("failed to scan request", slog.Any("err", err))
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *requestRepo) GetByID(id int64) (repository.Request, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, phone, email, contact_type, message, created_at, updated_at
		FROM requests
		WHERE id = ?
	`, id)

	var req repository.Request

	err := row.Scan(
		&req.ID,
		&req.Name,
		&req.Phone,
		&req.Email,
		&req.ContactType,
		&req.Message,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("request not found", slog.Int64("id", id))
			return req, repository.ErrNotFound
		}
		r.Logger.Error("failed to get request", slog.Int64("id", id), slog.Any("err", err))
		return repository.Request{}, err
	}

	return req, nil
}

func (r *requestRepo) Create(request repository.Request) (int64, error) {
	result, err := r.DB.Exec(`
		INSERT INTO requests(name, phone, email, contact_type, message, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		request.Name,
		request.Phone,
		request.Email,
		request.ContactType,
		request.Message,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		r.Logger.Error("failed to create request", slog.Any("err", err))
		return 0, err
	}

	r.Logger.Debug("request created", slog.String("user", request.Name))
	return result.LastInsertId()
}

func (r *requestRepo) Update(request repository.Request) error {
	_, err := r.DB.Exec(`
		UPDATE requests
		SET name = ?, phone = ?, email = ?, contact_type = ?, message = ?, updated_at = ?
		WHERE id = ?
	`,
		request.Name,
		request.Phone,
		request.Email,
		request.ContactType,
		request.Message,
		time.Now(),
		request.ID,
	)
	if err != nil {
		r.Logger.Error("failed to update request", slog.Int64("id", request.ID), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("request updated", slog.Int64("id", request.ID), slog.String("user", request.Name))
	return nil
}

func (r *requestRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM requests WHERE id = ?`, id)
	if err != nil {
		r.Logger.Error("failed to delete request", slog.Int64("id", id), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("request deleted", slog.Int64("id", id))
	return nil
}
