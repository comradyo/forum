package repository

import (
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
