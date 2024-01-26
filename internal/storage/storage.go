package storage

import (
	"database/sql"

	"github.com/modaniru/twitch-auth-server/internal/entity"
	"github.com/modaniru/twitch-auth-server/internal/storage/repo"
)

type User interface {
	CreateUser(userId, clientId string) (int, error)
	DeleteUserById(id int) error
	GetUserById(id int) (*entity.User, error)
	GetUserByUserId(userId string) (*entity.User, error)
}

type Storage struct {
	User
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{User: repo.NewUserStorage(db)}
}
