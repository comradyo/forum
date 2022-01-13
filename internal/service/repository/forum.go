package repository

import (
	"forum/forum/internal/models"
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

func (r *ForumRepository) CreateForum(forum *models.Forum) (*models.Forum, error) {
	err := r.db.QueryRow("QueryCreateForum", forum.Title, forum.User, forum.Slug).Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads,
	)
	if err != nil {
		return nil, models.ErrDatabase
	}
	return forum, nil
}

func (r *ForumRepository) GetForumDetails(slug string) (*models.Forum, error) {
	foundForum := &models.Forum{}
	err := r.db.QueryRow("QueryGetForumDetails", slug).Scan(&foundForum.Title, &foundForum.User, &foundForum.Slug, &foundForum.Posts, &foundForum.Threads)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrForumNotFound
		} else {
			return nil, models.ErrDatabase
		}
	}
	return foundForum, nil
}

func (r *ForumRepository) CreateForumThread(thread *models.Thread) (*models.Thread, error) {
	err := r.db.QueryRow("QueryCreateForumThread", thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)
	if err != nil {
		return nil, models.ErrDatabase
	}
	return thread, nil
}

func (r *ForumRepository) GetForumUsers(slug string, limit string, since string, desc string) (*models.Users, error) {
	users := &models.Users{}

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			rows, err = r.db.Query("QueryGetForumUsersSinceDesc", slug, since, limit)
		} else {
			rows, err = r.db.Query("QueryGetForumUsersDesc", slug, limit)
		}
	} else {
		if since != "" {
			rows, err = r.db.Query("QueryGetForumUsersSince", slug, since, limit)
		} else {
			rows, err = r.db.Query("QueryGetForumUsers", slug, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, models.ErrDatabase
	}

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			rows.Close()
			return nil, models.ErrDatabase
		}
		users.Users = append(users.Users, *user)
	}

	rows.Close()
	return users, nil
}

func (r *ForumRepository) GetForumThreads(slug string, limit string, since string, desc string) (*models.Threads, error) {
	threads := &models.Threads{}

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			rows, err = r.db.Query("QueryGetForumThreadsSinceDesc", slug, since, limit)
		} else {
			rows, err = r.db.Query("QueryGetForumThreadsDesc", slug, limit)
		}
	} else {
		if since != "" {
			rows, err = r.db.Query("QueryGetForumThreadsSince", slug, since, limit)
		} else {
			rows, err = r.db.Query("QueryGetForumThreads", slug, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, models.ErrDatabase
	}

	for rows.Next() {
		thread := &models.Thread{}
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			rows.Close()
			return nil, models.ErrDatabase
		}
		threads.Threads = append(threads.Threads, *thread)
	}

	rows.Close()
	return threads, nil
}
