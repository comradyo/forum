package usecase

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	"github.com/gofrs/uuid"
)

const forumLogMessage = "usecase:forum:"

type ForumUseCase struct {
	repository service.ForumRepositoryInterface
	userRepo   service.UserRepositoryInterface
	threadRepo service.ThreadRepositoryInterface
}

func NewForumUseCase(repository service.ForumRepositoryInterface, userRepo service.UserRepositoryInterface, threadRepo service.ThreadRepositoryInterface) *ForumUseCase {
	return &ForumUseCase{
		repository: repository,
		userRepo:   userRepo,
		threadRepo: threadRepo,
	}
}

func (u *ForumUseCase) CreateForum(forum *models.Forum) (*models.Forum, error) {
	user, err := u.userRepo.GetUserProfile(forum.User)
	if err != nil {
		return nil, err
	}
	forum.User = user.Nickname
	if forum.Slug != "" {
		oldForum, err := u.GetForumDetails(forum.Slug)
		if err == nil {
			return oldForum, models.ErrForumExists
		} else if err == models.ErrForumNotFound {
			return u.repository.CreateForum(forum)
		} else {
			return nil, err
		}
	} else {
		slug, _ := uuid.NewV4()
		forum.Slug = slug.String()
	}
	return u.repository.CreateForum(forum)
}

func (u *ForumUseCase) GetForumDetails(slug string) (*models.Forum, error) {
	return u.repository.GetForumDetails(slug)
}

func (u *ForumUseCase) CreateForumThread(slug string, thread *models.Thread) (*models.Thread, error) {
	forum, err := u.repository.GetForumDetails(slug)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserProfile(thread.Author)
	if err != nil {
		return nil, err
	}
	thread.Forum = forum.Slug
	thread.Author = user.Nickname
	if thread.Slug != "" {
		oldThread, err := u.threadRepo.GetThreadDetails(thread.Slug)
		if err == nil {
			return oldThread, models.ErrThreadExists
		} else if err == models.ErrThreadNotFound {
			return u.repository.CreateForumThread(thread)
		} else {
			return nil, err
		}
	}
	return u.repository.CreateForumThread(thread)
}

func (u *ForumUseCase) GetForumUsers(slug string, limit string, since string, desc string) (*models.Users, error) {
	_, err := u.repository.GetForumDetails(slug)
	if err != nil {
		return nil, err
	}
	if limit == "" {
		limit = "100"
	}
	return u.repository.GetForumUsers(slug, limit, since, desc)
}

func (u *ForumUseCase) GetForumThreads(slug string, limit string, since string, desc string) (*models.Threads, error) {
	_, err := u.repository.GetForumDetails(slug)
	if err != nil {
		return nil, err
	}
	if limit == "" {
		limit = "100"
	}
	return u.repository.GetForumThreads(slug, limit, since, desc)
}
