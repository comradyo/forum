package repository

import (
	"forum/forum/internal/models"
	"strconv"
	"strings"

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

func getSlugOrId(slugOrId string) (bool, int64, string) {
	var isId bool
	var idInt64 int64
	var slug string
	idInt, err := strconv.Atoi(slugOrId)
	if err != nil {
		isId = false
		slug = slugOrId
	} else {
		isId = true
		idInt64 = int64(idInt)
	}
	return isId, idInt64, slug
}

//TODO: Проверить
func (r *ThreadRepository) CreateThreadPosts(slugOrId string, posts *models.Posts) (*models.Posts, error) {
	isId, idInt64, slug := getSlugOrId(slugOrId)
	_ = isId
	_ = idInt64
	_ = slug

	return nil, nil
}

func (r *ThreadRepository) GetThreadDetails(slugOrId string) (*models.Thread, error) {
	isId, idInt64, slug := getSlugOrId(slugOrId)
	var row *pgx.Row

	if isId {
		query := `select * from thread where id = $1`
		row = r.db.QueryRow(query, idInt64)
	} else {
		query := `select * from thread where slug = $1`
		row = r.db.QueryRow(query, slug)
	}

	foundThread := &models.Thread{}
	err := row.Scan(
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
		if err == pgx.ErrNoRows {
			return nil, models.ErrThreadNotFound
		} else {
			return nil, models.ErrDatabase
		}
	}

	return foundThread, nil
}

//TODO: Проверить
func (r *ThreadRepository) UpdateThreadDetails(slugOrId string, thread *models.Thread) (*models.Thread, error) {
	isId, idInt64, slug := getSlugOrId(slugOrId)
	var row *pgx.Row

	if isId {
		query := `update thread set title = $1, message = $2 where id = $3 returning id, title, author, forum, message, votes, slug, created`
		row = r.db.QueryRow(query, thread.Title, thread.Message, idInt64)
	} else {
		query := `update thread set title = $1, message = $2 where slug = $3 returning id, title, author, forum, message, votes, slug, created`
		row = r.db.QueryRow(query, thread.Title, thread.Message, slug)
	}

	updatedThread := &models.Thread{}
	err := row.Scan(
		&updatedThread.Id,
		&updatedThread.Title,
		&updatedThread.Author,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.Votes,
		&updatedThread.Slug,
		&updatedThread.Created,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrThreadNotFound
		} else {
			return nil, models.ErrDatabase
		}
	}

	return nil, nil
}

//TODO: Проверить
func (r *ThreadRepository) GetThreadPosts(slugOrId string, limit string, since string, sort string, desc string) (*models.Posts, error) {
	isId, idInt64, slug := getSlugOrId(slugOrId)
	_ = isId
	_ = idInt64
	_ = slug

	return nil, nil
}

//TODO: Проверить
func (r *ThreadRepository) VoteForThread(slug string, vote *models.Vote) error {
	query := `insert into vote (thread, "user", voice) values $1, $2, $3`
	_, err := r.db.Exec(query, slug, vote.Nickname, vote.Voice)
	if strings.Contains(err.Error(), "duplicate") {
		query = `update vote set voice = $1 where thread = $2 and "user = $3"`
		_, err = r.db.Exec(query, vote.Voice, slug, vote.Nickname)
		if err != nil {
			return models.ErrDatabase
		}
	} else {
		return models.ErrDatabase
	}
	return nil
}
