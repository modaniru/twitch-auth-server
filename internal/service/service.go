package service

import (
	"github.com/modaniru/twitch-auth-server/internal/entity"
)

type User interface {
	GetUserInformation(id int) (*entity.UserResponse, error)
	CreateOrGetUser(token string) (int, error)
}

type AuthService interface {
	Auth(accessToken string) (string, error)
	Validate(jwt string) (int, error)
}

type Service struct {
	AuthService
	User
}

type Dependencies struct {
	AuthService
	User
}

func NewService(d Dependencies) *Service {
	return &Service{
		AuthService: d.AuthService,
		User:        d.User,
	}
}
