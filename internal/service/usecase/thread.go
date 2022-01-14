package usecase

import (
	"forum/internal/models"
	"forum/internal/service"
)

const threadLogMessage = "usecase:thread:"

type ThreadUseCase struct {
	repository service.ThreadRepositoryInterface
	userRepo   service.UserRepositoryInterface
	postRepo   service.PostRepositoryInterface
}

func NewThreadUseCase(repository service.ThreadRepositoryInterface, userRepo service.UserRepositoryInterface, postRepo service.PostRepositoryInterface) *ThreadUseCase {
	return &ThreadUseCase{
		repository: repository,
		userRepo:   userRepo,
		postRepo:   postRepo,
	}
}

func (u *ThreadUseCase) CreateThreadPosts(slugOrId string, posts *models.Posts) (*models.Posts, error) {
	thread, err := u.repository.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	if len(posts.Posts) == 0 {
		return &models.Posts{}, nil
	}
	return u.repository.CreateThreadPosts(thread.Id, thread.Forum, posts)
}

func (u *ThreadUseCase) GetThreadDetails(slugOrId string) (*models.Thread, error) {
	return u.repository.GetThreadDetails(slugOrId)
}

func (u *ThreadUseCase) UpdateThreadDetails(slugOrId string, thread *models.Thread) (*models.Thread, error) {
	oldThread, err := u.repository.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	thread.Id = oldThread.Id
	if thread.Title == "" {
		thread.Title = oldThread.Title
	}
	if thread.Message == "" {
		thread.Message = oldThread.Message
	}
	return u.repository.UpdateThreadDetails(thread.Id, thread)
}

func (u *ThreadUseCase) GetThreadPosts(slugOrId string, limit string, since string, sort string, desc string) (*models.Posts, error) {
	thread, err := u.repository.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	if limit == "" {
		limit = "100"
	}
	return u.repository.GetThreadPosts(thread.Id, limit, since, sort, desc)
}

func (u *ThreadUseCase) VoteForThread(slugOrId string, vote *models.Vote) (*models.Thread, error) {
	thread, err := u.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserProfile(vote.Nickname)
	if err != nil {
		return nil, err
	}
	vote.Nickname = user.Nickname
	err = u.repository.VoteForThread(thread.Id, vote)
	if err != nil {
		return nil, err
	}
	return u.repository.GetThreadDetails(slugOrId)
}
