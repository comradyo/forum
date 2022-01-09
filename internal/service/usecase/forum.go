package usecase

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
)

const forumLogMessage = "usecase:forum:"

type ForumUseCase struct {
	repository service.ForumRepositoryInterface
	userRepo   service.UserRepositoryInterface
}

func NewForumUseCase(repository service.ForumRepositoryInterface, userRepo service.UserRepositoryInterface) *ForumUseCase {
	return &ForumUseCase{
		repository: repository,
		userRepo:   userRepo,
	}
}

func (u *ForumUseCase) CreateForum(forum *models.Forum) (*models.Forum, error) {
	_, err := u.userRepo.GetUserProfile(forum.User)
	if err != nil {
		return nil, err
	}
	//TODO: Проверка, пустой ли слаг. Если пустой, то сгенерировать его
	return u.repository.CreateForum(forum)
}

func (u *ForumUseCase) GetForumDetails(slug string) (*models.Forum, error) {
	return u.repository.GetForumDetails(slug)
}

func (u *ForumUseCase) CreateForumThread(slug string, thread *models.Thread) (*models.Thread, error) {
	_, err := u.userRepo.GetUserProfile(thread.Author)
	if err != nil {
		return nil, err
	}
	_, err = u.repository.GetForumDetails(slug)
	if err != nil {
		return nil, err
	}
	thread.Forum = slug
	//TODO: Слаг для треда
	//thread.Slug = "123"
	return u.repository.CreateForumThread(thread)
}

func (u *ForumUseCase) GetForumUsers(slug string, limit string, since string, desc string) (*models.Users, error) {
	//TODO: мб лишние действия
	_, err := u.repository.GetForumDetails(slug)
	if err != nil {
		return nil, err
	}
	return u.repository.GetForumUsers(slug, limit, since, desc)
}

func (u *ForumUseCase) GetForumThreads(slug string, limit string, since string, desc string) (*models.Threads, error) {
	_, err := u.repository.GetForumDetails(slug)
	if err != nil {
		return nil, err
	}
	return u.repository.GetForumThreads(slug, limit, since, desc)
}
