package repository

import (
	"forum/forum/internal/models"
	"github.com/jackc/pgx"
)

const serviceLogMessage = "repository:service:"

type ServiceRepository struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (r *ServiceRepository) Clear() error {
	query := `truncate forum_user, vote, post, thread, forum, "user"`
	_, err := r.db.Exec(query)
	if err != nil {
		return models.ErrPostgres
	}
	return nil
}

func (r *ServiceRepository) GetStatus() (*models.Status, error) {
	status := &models.Status{}
	err := r.db.QueryRow(`select (*) count from "user"`).Scan(&status.User)
	if err != nil {
		return nil, models.ErrPostgres
	}
	err = r.db.QueryRow(`select (*) count from forum`).Scan(&status.Forum)
	if err != nil {
		return nil, models.ErrPostgres
	}
	err = r.db.QueryRow(`select (*) count from thread`).Scan(&status.Thread)
	if err != nil {
		return nil, models.ErrPostgres
	}
	err = r.db.QueryRow(`select (*) count from post`).Scan(&status.Post)
	if err != nil {
		return nil, models.ErrPostgres
	}
	return status, nil
}
