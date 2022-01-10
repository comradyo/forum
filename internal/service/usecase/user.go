package usecase

import (
	"forum/forum/internal/models"
	"forum/forum/internal/service"
)

const userLogMessage = "usecase:user:"

type UserUseCase struct {
	repository service.UserRepositoryInterface
}

func NewUserUseCase(repository service.UserRepositoryInterface) *UserUseCase {
	return &UserUseCase{
		repository: repository,
	}
}

func (u *UserUseCase) CreateUser(profile *models.User) (*models.User, error) {
	return u.repository.CreateUser(profile)
}

func (u *UserUseCase) GetUserProfile(nickname string) (*models.User, error) {
	return u.repository.GetUserProfile(nickname)
}

func (u *UserUseCase) UpdateUserProfile(profile *models.User) (*models.User, error) {
	oldProfile, err := u.GetUserProfile(profile.Nickname)
	if err != nil {
		return nil, err
	}
	if profile.Fullname == "" {
		profile.Fullname = oldProfile.Fullname
	}
	if profile.About == "" {
		profile.About = oldProfile.About
	}
	if profile.Email == "" {
		profile.Email = oldProfile.Email
	}
	return u.repository.UpdateUserProfile(profile)
}
