package repository

import (
	"database/sql"
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
	query := `insert into "forum" (title, user, slug) values ($1, $2, $3)`
	_, err := r.db.Exec(query, forum.Title, forum.User, forum.Slug)
	if err != nil {
		foundForum := &models.Forum{}
		query = `select (title, user, slug, posts, threads) from "forum" where slug = $1`
		err := r.db.QueryRow(query, forum.Slug).Scan(
			&foundForum.Title,
			&foundForum.User,
			&foundForum.Slug,
			&foundForum.Posts,
			&foundForum.Threads,
		)
		if err != nil {
			return nil, models.ErrPostgres
		}
		return foundForum, models.ErrForumExists
	}
	return forum, nil
}

func (r *ForumRepository) GetForumDetails(slug string) (*models.Forum, error) {
	query := `select (title, user, slug, posts, threads) from "forum" where slug = $1`
	rows, err := r.db.Query(query, slug)
	if err != nil {
		rows.Close()
		if err == sql.ErrNoRows {
			return nil, models.ErrForumNotFound
		} else {
			return nil, models.ErrPostgres
		}
	}
	foundForum := &models.Forum{}
	if rows.Next() {
		err = rows.Scan(
			&foundForum.Title,
			&foundForum.User,
			&foundForum.Slug,
			&foundForum.Posts,
			&foundForum.Threads,
		)
		if err != nil {
			rows.Close()
			return nil, models.ErrPostgres
		}
	}
	rows.Close()
	return foundForum, nil
}

func (r *ForumRepository) CreateForumThread(thread *models.Thread) (*models.Thread, error) {
	//TODO: Created либо автоматическим сделать, либо самому проставлять в usecase
	query := `insert into "thread" (title, author, forum, message, slug, created) values ($1, $2, $3, $4, $5, $6) returning id`
	rows, err := r.db.Query(query, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created)
	if err != nil {
		rows.Close()
		foundThread := &models.Thread{}
		query = `select * from "thread" where slug = $1`
		err := r.db.QueryRow(query, thread.Slug).Scan(
			&foundThread.Id,
			&foundThread.Title,
			&foundThread.Author,
			&foundThread.Forum,
			&foundThread.Message,
			&foundThread.Votes,
			&foundThread.Slug,
			&foundThread.Created,
		)
		if err != nil {
			return nil, models.ErrPostgres
		}
		return foundThread, models.ErrThreadExists
	}
	if rows.Next() {
		err = rows.Scan(&thread.Id)
		if err != nil {
			rows.Close()
			return nil, models.ErrPostgres
		}
	}
	rows.Close()
	return thread, nil
}

func (r *ForumRepository) GetForumUsers(slug string, limit string, since string, desc string) (*models.Users, error) {
	users := &models.Users{}
	query := `select (nickname, fullname, about, email) from "user"
				join forum_user on forum_user."user" = nickname
				where forum_user.forum = $1`

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query += ` and "user".nickname < $2 order by forum_user."user" desc`
			if limit != "" {
				query += ` limit $3`
				rows, err = r.db.Query(query, slug, since, limit)
			} else {
				rows, err = r.db.Query(query, slug, since)
			}
		} else {
			query += ` order by forum_user."user" desc`
			if limit != "" {
				query += ` limit $2`
				rows, err = r.db.Query(query, slug, limit)
			} else {
				rows, err = r.db.Query(query, slug)
			}
		}
	} else {
		if since != "" {
			query += ` and "user".nickname < $2 order by forum_user."user"`
			if limit != "" {
				query += ` limit $3`
				rows, err = r.db.Query(query, slug, since, limit)
			} else {
				rows, err = r.db.Query(query, slug, since)
			}
		} else {
			query += ` order by forum_user."user"`
			if limit != "" {
				query += ` limit $2`
				rows, err = r.db.Query(query, slug, limit)
			} else {
				rows, err = r.db.Query(query, slug)
			}
		}
	}

	if err != nil {
		rows.Close()
		if err == sql.ErrNoRows {
			return users, nil
		} else {
			return nil, models.ErrPostgres
		}
	}

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			rows.Close()
			return nil, models.ErrPostgres
		}
		users.Users = append(users.Users, *user)
	}

	rows.Close()
	return users, nil
}

func (r *ForumRepository) GetForumThreads(slug string, limit string, since string, desc string) (*models.Threads, error) {
	threads := &models.Threads{}
	query := `select (*) from thread where forum = $1`

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query += ` and created < $2 order by created desc`
			if limit != "" {
				query += ` limit $3`
				rows, err = r.db.Query(query, slug, since, limit)
			} else {
				rows, err = r.db.Query(query, slug, since)
			}
		} else {
			query += ` order by created desc`
			if limit != "" {
				query += ` limit $2`
				rows, err = r.db.Query(query, slug, limit)
			} else {
				rows, err = r.db.Query(query, slug)
			}
		}
	} else {
		if since != "" {
			query += ` and created < $2 order by created`
			if limit != "" {
				query += ` limit $3`
				rows, err = r.db.Query(query, slug, since, limit)
			} else {
				rows, err = r.db.Query(query, slug, since)
			}
		} else {
			query += ` order by created`
			if limit != "" {
				query += ` limit $2`
				rows, err = r.db.Query(query, slug, limit)
			} else {
				rows, err = r.db.Query(query, slug)
			}
		}
	}

	if err != nil {
		rows.Close()
		if err == sql.ErrNoRows {
			return threads, nil
		} else {
			return nil, models.ErrPostgres
		}
	}

	for rows.Next() {
		thread := &models.Thread{}
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			rows.Close()
			return nil, models.ErrPostgres
		}
		threads.Threads = append(threads.Threads, *thread)
	}

	rows.Close()
	return threads, nil
}
