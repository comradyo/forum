package repository

import (
	"github.com/jackc/pgx"
)

const forumLogMessage = "repository:forum:"

type ForumRepository struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}
