package repository

import (
	"forum/forum/internal/models"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"strings"
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
	query := `select nickname, fullname, about, email from "user" where nickname = $1 or email = $2`
	oldProfile := &models.User{}
	err := r.db.QueryRow(query, profile.Nickname, profile.Email).Scan(&oldProfile.Nickname, &oldProfile.Fullname, &oldProfile.About, &oldProfile.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			query := `insert into "user" (nickname, fullname, about, email) values ($1, $2, $3, $4)`
			_, err := r.db.Exec(query, profile.Nickname, profile.Fullname, profile.About, profile.Email)
			if err != nil {
				return nil, models.ErrPostgres
			}
			return profile, nil
		} else {
			return nil, models.ErrPostgres
		}
	} else {
		return oldProfile, models.ErrUserExists
	}
}

func (r *UserRepository) GetUserProfile(nickname string) (*models.User, error) {
	query := `select nickname, fullname, about, email from "user" where nickname = $1`
	foundProfile := &models.User{}
	err := r.db.QueryRow(query, nickname).Scan(&foundProfile.Nickname, &foundProfile.Fullname, &foundProfile.About, &foundProfile.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrUserNotFound
		} else {
			return nil, models.ErrPostgres
		}
	}
	return foundProfile, nil
}

func (r *UserRepository) UpdateUserProfile(profile *models.User) (*models.User, error) {
	query := `update "user" set fullname = $1, about = $2, email = $3 where nickname = $4 returning nickname, fullname, about, email`
	updatedProfile := &models.User{}
	err := r.db.QueryRow(query, profile.Fullname, profile.About, profile.Email, profile.Nickname).Scan(&updatedProfile.Nickname, &updatedProfile.Fullname, &updatedProfile.About, &updatedProfile.Email)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, models.ErrProfileUpdateConflict
		} else {
			return nil, models.ErrPostgres
		}
	}
	return updatedProfile, nil
}
