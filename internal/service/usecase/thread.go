package usecase

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
)

const threadLogMessage = "usecase:thread:"

type ThreadUseCase struct {
	repository service.ThreadRepositoryInterface
	userRepo   service.UserRepositoryInterface
}

func NewThreadUseCase(repository service.ThreadRepositoryInterface, userRepo service.UserRepositoryInterface) *ThreadUseCase {
	return &ThreadUseCase{
		repository: repository,
		userRepo:   userRepo,
	}
}

func (u *ThreadUseCase) CreateThreadPosts(slugOrId string, posts *models.Posts) (*models.Posts, error) {
	thread, err := u.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	for i, _ := range posts.Posts {
		posts.Posts[i].Forum = thread.Forum
		posts.Posts[i].Thread = thread.Id
	}
	//TODO: разобраться с parent. Скорее всего, это тоже делается в бд процедурой
	//TODO: posts.Created генерировать либо тут, либо в базе данных
	return u.repository.CreateThreadPosts(slugOrId, posts)
}

func (u *ThreadUseCase) GetThreadDetails(slugOrId string) (*models.Thread, error) {
	return u.repository.GetThreadDetails(slugOrId)
}

func (u *ThreadUseCase) UpdateThreadDetails(slugOrId string, thread *models.Thread) (*models.Thread, error) {
	//TODO: мб лишние действия
	oldThread, err := u.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	if thread.Title == "" {
		thread.Title = oldThread.Title
	}
	if thread.Message == "" {
		thread.Message = oldThread.Message
	}
	return u.repository.UpdateThreadDetails(slugOrId, thread)
}

func (u *ThreadUseCase) GetThreadPosts(slugOrId string, limit string, since string, sort string, desc string) (*models.Posts, error) {
	_, err := u.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	return u.repository.GetThreadPosts(slugOrId, limit, since, sort, desc)
}

func (u *ThreadUseCase) VoteForThread(slugOrId string, vote *models.Vote) (*models.Thread, error) {
	thread, err := u.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	_, err = u.userRepo.GetUserProfile(vote.Nickname)
	if err != nil {
		return nil, err
	}
	err = u.repository.VoteForThread(thread.Slug, vote)
	if err != nil {
		return nil, err
	}
	return u.GetThreadDetails(slugOrId)
}
