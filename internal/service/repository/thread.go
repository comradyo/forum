package repository

import (
	"github.com/jackc/pgx"
)

const threadLogMessage = "repository:thread:"

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}
