package repository

import (
	"forum/internal/models"

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

func (r *PostRepository) GetPost(id int64) (*models.Post, error) {
	foundPost := &models.Post{}
	err := r.db.QueryRow("QueryGetPost", id).Scan(
		&foundPost.Id,
		&foundPost.Parent,
		&foundPost.Author,
		&foundPost.Message,
		&foundPost.IsEdited,
		&foundPost.Forum,
		&foundPost.Thread,
		&foundPost.Created,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrPostNotFound
		} else {
			return nil, models.ErrDatabase
		}
	}
	return foundPost, nil
}

func (r *PostRepository) UpdatePostDetails(post *models.Post) (*models.Post, error) {
	updatedPost := &models.Post{}
	err := r.db.QueryRow("QueryUpdatePostDetails", post.Message, post.Id).Scan(
		&updatedPost.Id,
		&updatedPost.Parent,
		&updatedPost.Author,
		&updatedPost.Message,
		&updatedPost.IsEdited,
		&updatedPost.Forum,
		&updatedPost.Thread,
		&updatedPost.Created,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrPostNotFound
		} else {
			return nil, models.ErrDatabase
		}
	}
	return updatedPost, nil
}
