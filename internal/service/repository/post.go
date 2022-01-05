package repository

import (
	"github.com/jackc/pgx"
)

const postLogMessage = "repository:post:"

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) *PostRepository {
	return &PostRepository{
		db: db,
	}
}
