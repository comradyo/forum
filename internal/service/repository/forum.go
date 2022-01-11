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
	query := `insert into "forum" (title, "user", slug) values ($1, $2, $3) returning title, "user", slug, posts, threads`
	err := r.db.QueryRow(query, forum.Title, forum.User, forum.Slug).Scan(
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
	query := `select title, "user", slug, posts, threads from "forum" where slug = $1`
	foundForum := &models.Forum{}
	err := r.db.QueryRow(query, slug).Scan(&foundForum.Title, &foundForum.User, &foundForum.Slug, &foundForum.Posts, &foundForum.Threads)
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
	query := `insert into "thread" (title, author, forum, message, slug, created) values ($1, $2, $3, $4, $5, $6)
              returning id, title, author, forum, message, votes, slug, created`
	err := r.db.QueryRow(query, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).Scan(
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

//TODO: Проверить
func (r *ForumRepository) GetForumUsers(slug string, limit string, since string, desc string) (*models.Users, error) {
	users := &models.Users{}
	query := `select nickname, fullname, about, email from "user"
				join forum_user on forum_user."user" = nickname
				where forum_user.forum = $1`

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query += ` and "user".nickname < $2 order by forum_user."user" desc`
			query += ` limit $3`
			rows, err = r.db.Query(query, slug, since, limit)
		} else {
			query += ` order by forum_user."user" desc`
			query += ` limit $2`
			rows, err = r.db.Query(query, slug, limit)
		}
	} else {
		if since != "" {
			query += ` and "user".nickname > $2 order by forum_user."user"`
			query += ` limit $3`
			rows, err = r.db.Query(query, slug, since, limit)
		} else {
			query += ` order by forum_user."user"`
			query += ` limit $2`
			rows, err = r.db.Query(query, slug, limit)
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

//TODO: Проверить работу со временем
func (r *ForumRepository) GetForumThreads(slug string, limit string, since string, desc string) (*models.Threads, error) {
	threads := &models.Threads{}
	query := `select * from thread where forum = $1`

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query += ` and created <= $2 order by created desc`
			query += ` limit $3`
			rows, err = r.db.Query(query, slug, since, limit)
		} else {
			query += ` order by created desc`
			query += ` limit $2`
			rows, err = r.db.Query(query, slug, limit)
		}
	} else {
		if since != "" {
			query += ` and created >= $2 order by created`
			query += ` limit $3`
			rows, err = r.db.Query(query, slug, since, limit)
		} else {
			query += ` order by created`
			query += ` limit $2`
			rows, err = r.db.Query(query, slug, limit)
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
