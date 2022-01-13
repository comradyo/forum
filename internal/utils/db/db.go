package db

import (
	"github.com/jackc/pgx"
)

type queriesInfo struct {
	queryName string
	query     string
}

const (
	/////////////////////
	////////FORUM////////
	/////////////////////
	QueryCreateForum = `insert into "forum" (title, "user", slug) 
							values ($1, $2, $3) 
							returning title, "user", slug, posts, threads`
	/////////////////////
	QueryGetForumDetails = `select title, "user", slug, posts, threads from "forum" 
							where slug = $1`
	/////////////////////
	QueryCreateForumThread = `insert into "thread" (title, author, forum, message, slug, created) 
								values ($1, $2, $3, $4, $5, $6) 
								returning id, title, author, forum, message, votes, slug, created`
	/////////////////////
	QueryGetForumUsers = `select nickname, fullname, about, email from "user"
								join forum_user on "user".nickname = forum_user."user"
								and forum_user.forum = $1 
								order by nickname limit $2`
	/////////////////////
	QueryGetForumUsersSince = `select nickname, fullname, about, email from "user"
								join forum_user on "user".nickname = forum_user."user"
								and forum_user.forum = $1 
								and nickname > $2 
								order by nickname limit $3`
	/////////////////////
	QueryGetForumUsersDesc = `select nickname, fullname, about, email from "user"
								join forum_user on "user".nickname = forum_user."user"
								and forum_user.forum = $1 
								order by nickname desc limit $2`
	/////////////////////
	QueryGetForumUsersSinceDesc = `select nickname, fullname, about, email from "user"
									join forum_user on "user".nickname = forum_user."user"
									and forum_user.forum = $1 
									and nickname < $2 
									order by nickname desc limit $3`
	/////////////////////
	QueryGetForumThreads = `select * from thread 
							where forum = $1 
							order by created limit $2`
	/////////////////////
	QueryGetForumThreadsSince = `select * from thread 
									where forum = $1 
									and created >= $2 
									order by created limit $3`
	/////////////////////
	QueryGetForumThreadsDesc = `select * from thread 
								where forum = $1 
								order by created desc limit $2`
	/////////////////////
	QueryGetForumThreadsSinceDesc = `select * from thread 
										where forum = $1 
										and created <= $2 
										order by created desc limit $3`
	/////////////////////
	/////////POST////////
	/////////////////////
	QueryGetPost = `select id, parent, author, message, is_edited, forum, thread, created from post 
					where id = $1`
	/////////////////////
	QueryUpdatePostDetails = `update post set message = $1, is_edited = true 
								where id = $2 
								returning id, parent, author, message, is_edited, forum, thread, created`
	/////////////////////
	///////SERVICE///////
	/////////////////////
	QueryClear = `truncate forum_user, vote, post, thread, forum, "user"`
	/////////////////////
	QueryGetStatusUser   = `select count(*) from "user"`
	QueryGetStatusForum  = `select count(*) from forum`
	QueryGetStatusThread = `select count(*) from thread`
	QueryGetStatusPost   = `select count(*) from post`
	/////////////////////
	///////THREAD////////
	/////////////////////
	QueryUpdateThreadDetails = `update thread set title = $1, message = $2 
								where id = $3 
								returning id, title, author, forum, message, votes, slug, created`
	/////////////////////
	QueryGetThreadPostsFlat = `SELECT id, parent, author, message, is_edited, forum, thread, created
								from post 
								where thread = $1 
								order by id limit $2`
	/////////////////////
	QueryGetThreadPostsFlatSince = `SELECT id, parent, author, message, is_edited, forum, thread, created
									from post 
									where thread = $1 
									and id > $2 
									order by id limit $3`
	/////////////////////
	QueryGetThreadPostsFlatDesc = `SELECT id, parent, author, message, is_edited, forum, thread, created
									from post 
									where thread = $1 
									order by id desc limit $2`
	/////////////////////
	QueryGetThreadPostsFlatSinceDesc = `SELECT id, parent, author, message, is_edited, forum, thread, created 
										from post 
										where thread = $1 
										and id < $2
										order by id desc limit $3`
	/////////////////////
	QueryGetThreadPostsTree = `SELECT id, parent, author, message, is_edited, forum, thread, created 
								from post 
								where thread = $1 
								order by path, id limit $2`
	/////////////////////
	QueryGetThreadPostsTreeSince = `SELECT id, parent, author, message, is_edited, forum, thread, created 
									from post 
									where thread = $1 
									and path > (select path FROM post where id = $2) 
									order by path, id limit $3`
	/////////////////////
	QueryGetThreadPostsTreeDesc = `SELECT id, parent, author, message, is_edited, forum, thread, created 
									from post 
									where thread = $1 
									order by path desc, id desc limit $2`
	/////////////////////
	QueryGetThreadPostsTreeSinceDesc = `SELECT id, parent, author, message, is_edited, forum, thread, created 
										from post 
										where thread = $1 
										and path < (select path FROM post where id = $2) 
										order by path desc, id desc limit $3`
	/////////////////////
	QueryGetThreadPostsParentTree = `select id, parent, author, message, is_edited, forum, thread, created 
										from post 
										where path[1] in 
											(select id from post where thread = $1 and parent = 0 order by id limit $2)
										order by path, id`
	/////////////////////
	QueryGetThreadPostsParentTreeSince = `select id, parent, author, message, is_edited, forum, thread, created 
											from post 
											where path[1] in 
												(select id from post where thread = $1 and parent = 0 and path[1] > 
												(select path[1] from post where id = $2) 
												order by id limit $3)
											order by path, id`
	/////////////////////
	QueryGetThreadPostsParentTreeDesc = `select id, parent, author, message, is_edited, forum, thread, created 
											from post where path[1] in 
												(select id from post where thread = $1 and parent = 0 
												order by id desc limit $2)
											order by path[1] desc, path, id`
	/////////////////////
	QueryGetThreadPostsParentTreeSinceDesc = `select id, parent, author, message, is_edited, forum, thread, created 
												from post where path[1] in 
												(select id from post where thread = $1 and parent = 0 and path[1] <
												(select path[1] from post where id = $2) order by id desc limit $3)
												order by path[1] desc, path, id`
	/////////////////////
	QueryInsertVote = `insert into vote (thread, "user", voice) 
						values ($1, $2, $3)`
	/////////////////////
	QueryUpdateVote = `update vote set voice = $1 
						where thread = $2 
						and "user" = $3`
	/////////////////////
	////////USER/////////
	/////////////////////
	QueryCreateUser = `insert into "user" (nickname, fullname, about, email) 
						values ($1, $2, $3, $4)`
	/////////////////////
	QueryGetUserProfile = `select nickname, fullname, about, email from "user" 
							where nickname = $1`
	/////////////////////
	QueryGetUserProfileByMail = `select nickname, fullname, about, email from "user" 
									where email = $1`
	/////////////////////
	QueryUpdateUserProfile = `update "user" set fullname = $1, about = $2, email = $3 
								where nickname = $4 
								returning nickname, fullname, about, email`
)

var queries = []queriesInfo{
	{
		queryName: "QueryCreateForum",
		query:     QueryCreateForum,
	},
	{
		queryName: "QueryGetForumDetails",
		query:     QueryGetForumDetails,
	},
	{
		queryName: "QueryCreateForumThread",
		query:     QueryCreateForumThread,
	},
	{
		queryName: "QueryGetForumUsers",
		query:     QueryGetForumUsers,
	},
	{
		queryName: "QueryGetForumUsersSince",
		query:     QueryGetForumUsersSince,
	},
	{
		queryName: "QueryGetForumUsersDesc",
		query:     QueryGetForumUsersDesc,
	},
	{
		queryName: "QueryGetForumUsersSinceDesc",
		query:     QueryGetForumUsersSinceDesc,
	},
	{
		queryName: "QueryGetForumThreads",
		query:     QueryGetForumThreads,
	},
	{
		queryName: "QueryGetForumThreadsSince",
		query:     QueryGetForumThreadsSince,
	},
	{
		queryName: "QueryGetForumThreadsDesc",
		query:     QueryGetForumThreadsDesc,
	},
	{
		queryName: "QueryGetForumThreadsSinceDesc",
		query:     QueryGetForumThreadsSinceDesc,
	},
	{
		queryName: "QueryGetPost",
		query:     QueryGetPost,
	},
	{
		queryName: "QueryUpdatePostDetails",
		query:     QueryUpdatePostDetails,
	},
	{
		queryName: "QueryClear",
		query:     QueryClear,
	},
	{
		queryName: "QueryGetStatusUser",
		query:     QueryGetStatusUser,
	},
	{
		queryName: "QueryGetStatusForum",
		query:     QueryGetStatusForum,
	},
	{
		queryName: "QueryGetStatusThread",
		query:     QueryGetStatusThread,
	},
	{
		queryName: "QueryGetStatusPost",
		query:     QueryGetStatusPost,
	},
	{
		queryName: "QueryUpdateThreadDetails",
		query:     QueryUpdateThreadDetails,
	},
	{
		queryName: "QueryGetThreadPostsFlat",
		query:     QueryGetThreadPostsFlat,
	},
	{
		queryName: "QueryGetThreadPostsFlatSince",
		query:     QueryGetThreadPostsFlatSince,
	},
	{
		queryName: "QueryGetThreadPostsFlatDesc",
		query:     QueryGetThreadPostsFlatDesc,
	},
	{
		queryName: "QueryGetThreadPostsFlatSinceDesc",
		query:     QueryGetThreadPostsFlatSinceDesc,
	},
	{
		queryName: "QueryGetThreadPostsTree",
		query:     QueryGetThreadPostsTree,
	},
	{
		queryName: "QueryGetThreadPostsTreeSince",
		query:     QueryGetThreadPostsTreeSince,
	},
	{
		queryName: "QueryGetThreadPostsTreeDesc",
		query:     QueryGetThreadPostsTreeDesc,
	},
	{
		queryName: "QueryGetThreadPostsTreeSinceDesc",
		query:     QueryGetThreadPostsTreeSinceDesc,
	},
	{
		queryName: "QueryGetThreadPostsParentTree",
		query:     QueryGetThreadPostsParentTree,
	},
	{
		queryName: "QueryGetThreadPostsParentTreeSince",
		query:     QueryGetThreadPostsParentTreeSince,
	},
	{
		queryName: "QueryGetThreadPostsParentTreeDesc",
		query:     QueryGetThreadPostsParentTreeDesc,
	},
	{
		queryName: "QueryGetThreadPostsParentTreeSinceDesc",
		query:     QueryGetThreadPostsParentTreeSinceDesc,
	},
	{
		queryName: "QueryInsertVote",
		query:     QueryInsertVote,
	},
	{
		queryName: "QueryUpdateVote",
		query:     QueryUpdateVote,
	},
	{
		queryName: "QueryCreateUser",
		query:     QueryCreateUser,
	},
	{
		queryName: "QueryGetUserProfile",
		query:     QueryGetUserProfile,
	},
	{
		queryName: "QueryGetUserProfileByMail",
		query:     QueryGetUserProfileByMail,
	},
	{
		queryName: "QueryUpdateUserProfile",
		query:     QueryUpdateUserProfile,
	},
}

func NewConnPool() (*pgx.ConnPool, error) {
	config := pgx.ConnConfig{
		User:                 "postgres",
		Database:             "postgres",
		Password:             "password",
		PreferSimpleProtocol: false,
	}
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	return pgx.NewConnPool(connPoolConfig)
}

func Prepare(db *pgx.ConnPool) error {
	for _, query := range queries {
		_, err := db.Prepare(query.queryName, query.query)
		if err != nil {
			return err
		}
	}
	return nil
}
