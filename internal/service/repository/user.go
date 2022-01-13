package repository

import (
	"forum/forum/internal/models"
	"strings"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
)

const userLogMessage = "repository:user:"

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(profile *models.User) (*models.User, error) {
	_, err := r.db.Exec("QueryCreateUser", profile.Nickname, profile.Fullname, profile.About, profile.Email)
	if err != nil {
		return nil, models.ErrDatabase
	}
	return profile, nil
}

func (r *UserRepository) GetUserProfile(nickname string) (*models.User, error) {
	foundProfile := &models.User{}
	err := r.db.QueryRow("QueryGetUserProfile", nickname).Scan(&foundProfile.Nickname, &foundProfile.Fullname, &foundProfile.About, &foundProfile.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrUserNotFound
		} else {
			return nil, models.ErrDatabase
		}
	}
	return foundProfile, nil
}

func (r *UserRepository) GetUserProfileByMail(email string) (*models.User, error) {
	foundProfile := &models.User{}
	err := r.db.QueryRow("QueryGetUserProfileByMail", email).Scan(&foundProfile.Nickname, &foundProfile.Fullname, &foundProfile.About, &foundProfile.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrUserNotFound
		} else {
			return nil, models.ErrDatabase
		}
	}
	return foundProfile, nil
}

func (r *UserRepository) UpdateUserProfile(profile *models.User) (*models.User, error) {
	updatedProfile := &models.User{}
	err := r.db.QueryRow("QueryUpdateUserProfile", profile.Fullname, profile.About, profile.Email, profile.Nickname).Scan(&updatedProfile.Nickname, &updatedProfile.Fullname, &updatedProfile.About, &updatedProfile.Email)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, models.ErrProfileUpdateConflict
		} else {
			return nil, models.ErrDatabase
		}
	}
	return updatedProfile, nil
}
