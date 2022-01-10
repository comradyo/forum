package repository

import (
	"forum/forum/internal/models"

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

//TODO: Проверить
func (r *PostRepository) GetPost(id int64) (*models.Post, error) {
	query := `select id, parent, author, message, is_edited, forum, thread, created from post where id = $1`
	foundPost := &models.Post{}
	err := r.db.QueryRow(query, id).Scan(
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

//TODO: Проверить
func (r *PostRepository) UpdatePostDetails(post *models.Post) (*models.Post, error) {
	query := `update post set message = $1, is_edited = true where id = $2 returning
				id, parent, author, message, is_edited, forum, thread, created`
	updatedPost := &models.Post{}
	err := r.db.QueryRow(query, post.Message, post.Id).Scan(
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
