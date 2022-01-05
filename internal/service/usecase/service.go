package usecase

import "forum/forum/internal/service"

const serviceLogMessage = "usecase:service:"

type ServiceUseCase struct {
	repository service.ServiceRepositoryInterface
}

func NewServiceUseCase(repository service.ServiceRepositoryInterface) *ServiceUseCase {
	return &ServiceUseCase{
		repository: repository,
	}
}
