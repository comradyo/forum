package repository

import (
	"github.com/jackc/pgx"
)

const userLogMessage = "repository:user:"

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
