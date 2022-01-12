package usecase

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	"time"
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
	thread, err := u.GetThreadDetails(slugOrId)
	if err != nil {
		return nil, err
	}
	if posts == nil {
		return nil, nil
	}
	created := time.Now()
	for i, _ := range posts.Posts {
		if posts.Posts[i].Parent != 0 {
			parentPost, err := u.postRepo.GetPost(posts.Posts[i].Parent)
			if err != nil {
				return nil, err
			}
			if parentPost.Thread != thread.Id {
				return nil, models.ErrPostNotFound
			}
		}
		posts.Posts[i].Forum = thread.Forum
		posts.Posts[i].Thread = thread.Id
		posts.Posts[i].Created = created
	}
	return u.repository.CreateThreadPosts(thread.Id, posts)
}

func (u *ThreadUseCase) GetThreadDetails(slugOrId string) (*models.Thread, error) {
	return u.repository.GetThreadDetails(slugOrId)
}

func (u *ThreadUseCase) UpdateThreadDetails(slugOrId string, thread *models.Thread) (*models.Thread, error) {
	oldThread, err := u.GetThreadDetails(slugOrId)
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
	thread, err := u.GetThreadDetails(slugOrId)
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
	_, err = u.userRepo.GetUserProfile(vote.Nickname)
	if err != nil {
		return nil, err
	}
	err = u.repository.VoteForThread(thread.Id, vote)
	if err != nil {
		return nil, err
	}
	return u.GetThreadDetails(thread.Slug)
}
