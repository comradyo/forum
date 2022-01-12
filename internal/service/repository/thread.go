package repository

import (
	"fmt"
	"forum/forum/internal/models"
	log "forum/forum/pkg/logger"
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

func getSlugOrId(slugOrId string) (bool, int32, string) {
	var isId bool
	var id int32
	var slug string
	idInt, err := strconv.Atoi(slugOrId)
	if err != nil {
		isId = false
		slug = slugOrId
	} else {
		isId = true
		id = int32(idInt)
	}
	return isId, id, slug
}

func (r *ThreadRepository) CreateThreadPosts(id int32, posts *models.Posts) (*models.Posts, error) {
	createdPosts := &models.Posts{}

	query := `insert into post (parent, author, message, forum, thread, created) values`
	var values []interface{}

	for i, post := range posts.Posts {
		valuesNumbers := fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d),",
			i*6+1,
			i*6+2,
			i*6+3,
			i*6+4,
			i*6+5,
			i*6+6,
		)
		query += valuesNumbers
		values = append(values, post.Parent, post.Author, post.Message, post.Forum, post.Thread, post.Created)
	}
	query = strings.TrimSuffix(query, ",")
	query += ` returning id, parent, author, message, is_edited, forum, thread, created`

	rows, err := r.db.Query(query, values...)
	if err != nil {
		rows.Close()
		return nil, models.ErrDatabase
	}
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			rows.Close()
			return nil, models.ErrDatabase
		}
		createdPosts.Posts = append(createdPosts.Posts, *post)
	}
	rows.Close()
	return createdPosts, nil
}

func (r *ThreadRepository) GetThreadDetails(slugOrId string) (*models.Thread, error) {
	isId, id, slug := getSlugOrId(slugOrId)
	var row *pgx.Row

	if isId {
		query := `select * from thread where id = $1`
		row = r.db.QueryRow(query, id)
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

func (r *ThreadRepository) UpdateThreadDetails(id int32, thread *models.Thread) (*models.Thread, error) {
	query := `update thread set title = $1, message = $2 where id = $3 returning id, title, author, forum, message, votes, slug, created`
	row := r.db.QueryRow(query, thread.Title, thread.Message, id)
	err := row.Scan(
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

func (r *ThreadRepository) getThreadPostsFlat(id int32, limit string, since string, desc string) (*models.Posts, error) {
	posts := &models.Posts{}
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created 
				from post where thread = $1`

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query += ` and id < $2 order by id desc`
			query += ` limit $3`
			rows, err = r.db.Query(query, id, since, limit)
		} else {
			query += ` order by id desc`
			query += ` limit $2`
			rows, err = r.db.Query(query, id, limit)
		}
	} else {
		if since != "" {
			query += ` and id > $2 order by id`
			query += ` limit $3`
			rows, err = r.db.Query(query, id, since, limit)
		} else {
			query += ` order by id`
			query += ` limit $2`
			rows, err = r.db.Query(query, id, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, models.ErrDatabase
	}

	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			rows.Close()
			return nil, models.ErrDatabase
		}
		posts.Posts = append(posts.Posts, *post)
	}

	rows.Close()
	return posts, nil
}

func (r *ThreadRepository) getThreadPostsTree(id int32, limit string, since string, desc string) (*models.Posts, error) {
	posts := &models.Posts{}
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created 
				from post where thread = $1`

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query += ` and path < (select path FROM post where id = $2) order by path desc, id desc`
			query += ` limit $3`
			rows, err = r.db.Query(query, id, since, limit)
		} else {
			query += ` order by path desc, id desc`
			query += ` limit $2`
			rows, err = r.db.Query(query, id, limit)
		}
	} else {
		if since != "" {
			query += ` and path > (select path FROM post where id = $2) order by path, id`
			query += ` limit $3`
			rows, err = r.db.Query(query, id, since, limit)
		} else {
			query += ` order by path, id`
			query += ` limit $2`
			rows, err = r.db.Query(query, id, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, models.ErrDatabase
	}

	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			rows.Close()
			return nil, models.ErrDatabase
		}
		posts.Posts = append(posts.Posts, *post)
	}

	rows.Close()
	return posts, nil
}

func (r *ThreadRepository) getThreadPostsParentTree(id int32, limit string, since string, desc string) (*models.Posts, error) {
	posts := &models.Posts{}
	query := `select id, parent, author, message, is_edited, forum, thread, created 
				from post where path[1] in`

	var rows *pgx.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query += ` (select id from post where thread = $1 and parent = 0 and path[1] < 
						(select path[1] from post where id = $2) order by id desc limit $3)
						order by path[1] desc, path, id`
			rows, err = r.db.Query(query, id, since, limit)
		} else {
			query += ` (select id from post where thread = $1 and parent = 0 order by id desc limit $2)
						order by path[1] desc, path, id`
			rows, err = r.db.Query(query, id, limit)
		}
	} else {
		if since != "" {
			query += ` (select id from post where thread = $1 and parent = 0 and path[1] > 
						(select path[1] from post where id = $2) order by id limit $3)
						order by path, id`
			rows, err = r.db.Query(query, id, since, limit)
		} else {
			query += ` (select id from post where thread = $1 and parent = 0 order by id limit $2)
						order by path, id`
			rows, err = r.db.Query(query, id, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, models.ErrDatabase
	}

	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			rows.Close()
			return nil, models.ErrDatabase
		}
		posts.Posts = append(posts.Posts, *post)
	}

	rows.Close()
	return posts, nil
}

//TODO: Проверить
func (r *ThreadRepository) GetThreadPosts(id int32, limit string, since string, sort string, desc string) (*models.Posts, error) {
	if sort == "" || sort == "flat" {
		return r.getThreadPostsFlat(id, limit, since, desc)
	} else if sort == "tree" {
		return r.getThreadPostsTree(id, limit, since, desc)
	} else if sort == "parent_tree" {
		return r.getThreadPostsParentTree(id, limit, since, desc)
	} else {
		return nil, models.ErrDatabase
	}
}

/*
`insert into vote (thread, "user", voice) values ($1, $2, $3)`
`update vote set voice = $1 where thread = $2 and "user" = $3`
*/

func (r *ThreadRepository) VoteForThread(id int32, vote *models.Vote) error {
	query := `insert into vote (thread, "user", voice) values ($1, $2, $3)`
	_, err := r.db.Exec(query, id, vote.Nickname, vote.Voice)
	if err != nil {
		log.Debug("err 1 =", err)
		if strings.Contains(err.Error(), "duplicate") {
			query = `update vote set voice = $1 where thread = $2 and "user" = $3`
			_, err = r.db.Exec(query, vote.Voice, id, vote.Nickname)
			if err != nil {
				log.Debug("err 2 =", err)
				return models.ErrDatabase
			}
		} else {
			return models.ErrDatabase
		}
	}
	return nil
}
