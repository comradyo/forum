package repository

import (
	"database/sql"
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

func (r *PostRepository) GetPost(id int64) (*models.Post, error) {
	query := `select id, parent, author, message, is_edited, foru, thread, created from post where id = $1`
	rows, err := r.db.Query(query, id)
	if err != nil {
		rows.Close()
		if err == sql.ErrNoRows {
			return nil, models.ErrPostNotFound
		} else {
			return nil, models.ErrPostgres
		}
	}
	foundPost := &models.Post{}
	if rows.Next() {
		err = rows.Scan(
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
			rows.Close()
			return nil, models.ErrPostgres
		}
	}
	rows.Close()
	return foundPost, nil
}

/*
   id          serial primary key,
   parent      int default 0,
   author      citext references "user"(nickname) on delete cascade not null,
   message     text not null,
   is_edited   bool not null default false,
   forum       citext references "forum"(slug) on delete cascade not null,
   thread      int references "thread"(id) on delete cascade not null,
   created     timestamp with time zone default now(),
*/

func (r *PostRepository) UpdatePostDetails(post *models.Post) (*models.Post, error) {
	query := `update post set message = $1, is_edited = true where id = $2 returning
				id, parent, author, message, is_edited, forum, thread, created`
	rows, err := r.db.Query(query, post.Message, post.Id)
	if err != nil {
		rows.Close()
		if err == sql.ErrNoRows {
			return nil, models.ErrPostNotFound
		} else {
			return nil, models.ErrPostgres
		}
	}
	updatedPost := &models.Post{}
	if rows.Next() {
		err = rows.Scan(
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
			rows.Close()
			return nil, models.ErrPostgres
		}
	}
	rows.Close()
	return updatedPost, nil
}
