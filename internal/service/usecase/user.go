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

func (u *UserUseCase) CreateUser(profile *models.User) ([]models.User, error) {
	var users []models.User
	oldUser1, err1 := u.repository.GetUserProfile(profile.Nickname)
	if err1 == nil {
		users = append(users, *oldUser1)
	}
	oldUser2, err2 := u.repository.GetUserProfileByMail(profile.Email)
	if err2 == nil {
		if oldUser1 != nil {
			if oldUser1.Nickname != oldUser2.Nickname && oldUser1.Email != oldUser2.Email {
				users = append(users, *oldUser2)
			}
		} else {
			users = append(users, *oldUser2)
		}
	}
	if err1 == nil || err2 == nil {
		return users, models.ErrUserExists
	}
	newUser, err := u.repository.CreateUser(profile)
	if err != nil {
		return nil, err
	}
	users = append(users, *newUser)
	return users, nil
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
