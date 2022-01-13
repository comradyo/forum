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
	_, err := r.db.Exec("QueryClear")
	if err != nil {
		return models.ErrDatabase
	}
	return nil
}

func (r *ServiceRepository) GetStatus() (*models.Status, error) {
	status := &models.Status{}
	err := r.db.QueryRow("QueryGetStatusUser").Scan(&status.User)
	if err != nil {
		return nil, models.ErrDatabase
	}
	err = r.db.QueryRow("QueryGetStatusForum").Scan(&status.Forum)
	if err != nil {
		return nil, models.ErrDatabase
	}
	err = r.db.QueryRow("QueryGetStatusThread").Scan(&status.Thread)
	if err != nil {
		return nil, models.ErrDatabase
	}
	err = r.db.QueryRow("QueryGetStatusPost").Scan(&status.Post)
	if err != nil {
		return nil, models.ErrDatabase
	}
	return status, nil
}
