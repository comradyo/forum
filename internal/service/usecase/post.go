package usecase

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
	"strconv"
	"strings"
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
	post, err := u.repository.GetPost(id)
	if err != nil {
		return nil, err
	}
	postFull := &models.PostFull{Post: post}
	if strings.Contains(related, "user") {
		user, _ := u.userRepo.GetUserProfile(post.Author)
		postFull.Author = user
	}
	if strings.Contains(related, "forum") {
		forum, _ := u.forumRepo.GetForumDetails(post.Forum)
		postFull.Forum = forum
	}
	if strings.Contains(related, "thread") {
		threadIdStr := strconv.Itoa(int(post.Thread))
		thread, _ := u.threadRepo.GetThreadDetails(threadIdStr)
		postFull.Thread = thread
	}
	return postFull, nil
}

func (u *PostUseCase) UpdatePostDetails(post *models.Post) (*models.Post, error) {
	postFull, err := u.GetPostDetails(post.Id, "")
	if err != nil {
		return nil, err
	}
	if post.Message != postFull.Post.Message {
		post.IsEdited = true
	} else {
		return postFull.Post, nil
	}
	return u.repository.UpdatePostDetails(post)
}
