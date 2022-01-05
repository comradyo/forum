package usecase

import "forum/forum/internal/service"

const postLogMessage = "usecase:post:"

type PostUseCase struct {
	repository service.PostRepositoryInterface
}

func NewPostUseCase(repository service.PostRepositoryInterface) *PostUseCase {
	return &PostUseCase{
		repository: repository,
	}
}
