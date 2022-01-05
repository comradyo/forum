package usecase

import "forum/forum/internal/service"

const threadLogMessage = "usecase:thread:"

type ThreadUseCase struct {
	repository service.ThreadRepositoryInterface
}

func NewThreadUseCase(repository service.ThreadRepositoryInterface) *ThreadUseCase {
	return &ThreadUseCase{
		repository: repository,
	}
}
