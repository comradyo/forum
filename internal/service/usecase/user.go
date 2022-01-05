package usecase

import "forum/forum/internal/service"

const userLogMessage = "usecase:user:"

type UserUseCase struct {
	repository service.UserRepositoryInterface
}

func NewUserUseCase(repository service.UserRepositoryInterface) *UserUseCase {
	return &UserUseCase{
		repository: repository,
	}
}
