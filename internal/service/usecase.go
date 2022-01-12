package service

import "forum/forum/internal/models"

type ForumUseCaseInterface interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	GetForumDetails(slug string) (*models.Forum, error)
	CreateForumThread(slug string, thread *models.Thread) (*models.Thread, error)
	GetForumUsers(slug string, limit string, since string, desc string) (*models.Users, error)
	GetForumThreads(slug string, limit string, since string, desc string) (*models.Threads, error)
}

type PostUseCaseInterface interface {
	GetPostDetails(id int64, related string) (*models.PostFull, error)
	UpdatePostDetails(post *models.Post) (*models.Post, error)
}

type ServiceUseCaseInterface interface {
	Clear() error
	GetStatus() (*models.Status, error)
}

type ThreadUseCaseInterface interface {
	CreateThreadPosts(slugOrId string, posts *models.Posts) (*models.Posts, error)
	GetThreadDetails(slugOrId string) (*models.Thread, error)
	UpdateThreadDetails(slugOrId string, thread *models.Thread) (*models.Thread, error)
	GetThreadPosts(slugOrId string, limit string, since string, sort string, desc string) (*models.Posts, error)
	VoteForThread(slugOrId string, vote *models.Vote) (*models.Thread, error)
}

type UserUseCaseInterface interface {
	CreateUser(profile *models.User) ([]models.User, error)
	GetUserProfile(nickname string) (*models.User, error)
	UpdateUserProfile(profile *models.User) (*models.User, error)
}
