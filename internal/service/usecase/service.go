package usecase

import (
	"forum/internal/models"
	"forum/internal/service"
)

const serviceLogMessage = "usecase:service:"

type ServiceUseCase struct {
	repository service.ServiceRepositoryInterface
}

func NewServiceUseCase(repository service.ServiceRepositoryInterface) *ServiceUseCase {
	return &ServiceUseCase{
		repository: repository,
	}
}

func (u *ServiceUseCase) Clear() error {
	return u.repository.Clear()
}

func (u *ServiceUseCase) GetStatus() (*models.Status, error) {
	return u.repository.GetStatus()
}
