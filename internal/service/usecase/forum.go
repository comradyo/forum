package usecase

import "forum/forum/internal/service"

const forumLogMessage = "usecase:forum:"

type ForumUseCase struct {
	repository service.ForumRepositoryInterface
}

func NewForumUseCase(repository service.ForumRepositoryInterface) *ForumUseCase {
	return &ForumUseCase{
		repository: repository,
	}
}
