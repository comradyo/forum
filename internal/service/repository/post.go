package repository

import (
	"forum/internal/models"
	"strings"

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

func (r *PostRepository) GetPostFull(id int64, related string) (*models.PostFull, error) {
	postFull := &models.PostFull{}
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
	postFull.Post = foundPost
	if strings.Contains(related, "user") {
		foundProfile := &models.User{}
		_ = r.db.QueryRow("QueryGetUserProfile", foundPost.Author).Scan(&foundProfile.Nickname, &foundProfile.Fullname, &foundProfile.About, &foundProfile.Email)
		postFull.Author = foundProfile
	}
	if strings.Contains(related, "forum") {
		foundForum := &models.Forum{}
		_ = r.db.QueryRow("QueryGetForumDetails", foundPost.Forum).Scan(&foundForum.Title, &foundForum.User, &foundForum.Slug, &foundForum.Posts, &foundForum.Threads)
		postFull.Forum = foundForum
	}
	if strings.Contains(related, "thread") {
		foundThread := &models.Thread{}
		_ = r.db.QueryRow(`select * from thread where id = $1`, foundPost.Thread).Scan(
			&foundThread.Id,
			&foundThread.Title,
			&foundThread.Author,
			&foundThread.Forum,
			&foundThread.Message,
			&foundThread.Votes,
			&foundThread.Slug,
			&foundThread.Created,
		)
		postFull.Thread = foundThread
	}
	return postFull, nil
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
