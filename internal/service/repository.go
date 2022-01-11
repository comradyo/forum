package service

import "forum/forum/internal/models"

type ForumRepositoryInterface interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	GetForumDetails(slug string) (*models.Forum, error)
	CreateForumThread(thread *models.Thread) (*models.Thread, error)
	GetForumUsers(slug string, limit string, since string, desc string) (*models.Users, error)
	GetForumThreads(slug string, limit string, since string, desc string) (*models.Threads, error)
}

type PostRepositoryInterface interface {
	GetPost(id int64) (*models.Post, error)
	UpdatePostDetails(post *models.Post) (*models.Post, error)
}

type ServiceRepositoryInterface interface {
	Clear() error
	GetStatus() (*models.Status, error)
}

type ThreadRepositoryInterface interface {
	CreateThreadPosts(id int32, posts *models.Posts) (*models.Posts, error)
	GetThreadDetails(slugOrId string) (*models.Thread, error)
	UpdateThreadDetails(id int32, thread *models.Thread) (*models.Thread, error)
	GetThreadPosts(id int32, limit string, since string, sort string, desc string) (*models.Posts, error)
	VoteForThread(id int32, vote *models.Vote) error
}

type UserRepositoryInterface interface {
	CreateUser(profile *models.User) (*models.User, error)
	GetUserProfile(nickname string) (*models.User, error)
	UpdateUserProfile(profile *models.User) (*models.User, error)
}
