package repository

import (
	"context"
	"database/sql"

	"github.com/modaniru/twitch-auth-server/internal/db"
)

type UserRepositorier interface {
	CreateUser(params *db.CreateUserParams) (int, error)
	DeleteUser(id int) error
	GetUser(id int) (*db.User, error)
	GetUserByUserId(userId string) (*db.User, error)
}

type UserRepository struct {
	queries *db.Queries
	db      *sql.DB // for trunc
}

func NewUserRepository(db *sql.DB, queries *db.Queries) UserRepositorier {
	return &UserRepository{db: db, queries: queries}
}

func (u *UserRepository) CreateUser(params *db.CreateUserParams) (int, error) {
	id, err := u.queries.CreateUser(context.Background(), *params)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (u *UserRepository) DeleteUser(id int) error {
	return u.queries.DeleteUser(context.Background(), int32(id))
}

func (u *UserRepository) GetUser(id int) (*db.User, error) {
	user, err := u.queries.GetUser(context.Background(), int32(id))
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) GetUserByUserId(userId string) (*db.User, error) {
	user, err := u.queries.GetUserByUserId(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
