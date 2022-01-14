package usecase

import (
	"forum/internal/models"
	"forum/internal/service"
)

const postLogMessage = "usecase:post:"

type PostUseCase struct {
	repository service.PostRepositoryInterface
	userRepo   service.UserRepositoryInterface
	forumRepo  service.ForumRepositoryInterface
	threadRepo service.ThreadRepositoryInterface
}

func NewPostUseCase(
	repository service.PostRepositoryInterface,
	userRepo service.UserRepositoryInterface,
	forumRepo service.ForumRepositoryInterface,
	threadRepo service.ThreadRepositoryInterface,
) *PostUseCase {
	return &PostUseCase{
		repository: repository,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
		threadRepo: threadRepo,
	}
}

func (u *PostUseCase) GetPostDetails(id int64, related string) (*models.PostFull, error) {
	return u.repository.GetPostFull(id, related)
}

func (u *PostUseCase) UpdatePostDetails(post *models.Post) (*models.Post, error) {
	postFull, err := u.GetPostDetails(post.Id, "")
	if err != nil {
		return nil, err
	}
	if post.Message != postFull.Post.Message && post.Message != "" {
		post.IsEdited = true
	} else {
		return postFull.Post, nil
	}
	return u.repository.UpdatePostDetails(post)
}
