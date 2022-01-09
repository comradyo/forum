package repository

import (
	"database/sql"
	"forum/forum/internal/models"
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
	query := `insert into "user" (nickname, fullname, about, email) values ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, profile.Nickname, profile.Fullname, profile.About, profile.Email)
	if err != nil {
		foundProfile := &models.User{}
		query = `select (nickname, fullname, about, email) from "user" where nickname = $1 or email = $2`
		err := r.db.QueryRow(query, profile.Nickname, profile.Email).Scan(
			&foundProfile.Nickname,
			&foundProfile.Fullname,
			&foundProfile.About,
			&foundProfile.Email,
		)
		if err != nil {
			return nil, models.ErrPostgres
		}
		return foundProfile, models.ErrUserExists
	}
	return profile, nil
}

func (r *UserRepository) GetUserProfile(nickname string) (*models.User, error) {
	query := `select (nickname, fullname, about, email) from "user" where nickname = $1`
	rows, err := r.db.Query(query, nickname)
	if err != nil {
		rows.Close()
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		} else {
			return nil, models.ErrPostgres
		}
	}
	foundProfile := &models.User{}
	if rows.Next() {
		err = rows.Scan(&foundProfile.Nickname, &foundProfile.Fullname, &foundProfile.About, &foundProfile.Email)
		if err != nil {
			rows.Close()
			return nil, models.ErrPostgres
		}
	}
	rows.Close()
	return foundProfile, nil
}

func (r *UserRepository) UpdateUserProfile(profile *models.User) (*models.User, error) {
	query := `update "user" set fullname = $1, about = $2, email = $3 where nickname = $4 returning nickname, fullname, about, email`
	rows, err := r.db.Query(query, profile.Fullname, profile.About, profile.Email, profile.Nickname)
	if err != nil {
		rows.Close()
		//TODO: ErrUpdateConflict
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		} else {
			return nil, models.ErrPostgres
		}
	}
	updatedProfile := &models.User{}
	if rows.Next() {
		err = rows.Scan(&updatedProfile.Nickname, &updatedProfile.Fullname, &updatedProfile.About, &updatedProfile.Email)
		if err != nil {
			rows.Close()
			return nil, models.ErrPostgres
		}
	}
	rows.Close()
	return updatedProfile, nil
}
