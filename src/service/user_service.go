package service

import (
	"database/sql"

	"github.com/modaniru/twitch-auth-server/src/client"
	"github.com/modaniru/twitch-auth-server/src/db"
	"github.com/modaniru/twitch-auth-server/src/dto/response"
	"github.com/modaniru/twitch-auth-server/src/repository"
)

type UserServicer interface {
	GetUserInformation(id int) (*response.UserResponse, error)
	CreateOrGetUser(token string) (int, error)
}

type UserService struct {
	userRepository repository.UserRepositorier
	twitchClient   client.TwitchClienter
}

// TODO construct
func NewUserService(userRepository repository.UserRepositorier, twitchClient client.TwitchClienter) UserServicer {
	return &UserService{
		userRepository: userRepository,
		twitchClient:   twitchClient,
	}
}

func (u *UserService) CreateOrGetUser(token string) (int, error) {
	creds, err := u.twitchClient.Validate(token)
	if err != nil {
		return 0, err
	}
	params := db.CreateUserParams{
		UserID:   creds.UserId,
		ClientID: creds.ClientId,
		RegDate:  "test",
	}
	user, err := u.userRepository.GetUserByUserId(params.UserID)
	if err == sql.ErrNoRows {
		return u.userRepository.CreateUser(&params)
	}
	return int(user.ID), nil
}

func (u *UserService) GetUserInformation(id int) (*response.UserResponse, error) {
	creds, err := u.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}
	userInfo, err := u.twitchClient.GetUserInformationApp(&client.ValidateResponse{
		ClientId: creds.ClientID,
		UserId:   creds.UserID,
	})
	if err != nil {
		return nil, err
	}
	userColor, err := u.twitchClient.GetUserChatColorApp(&client.ValidateResponse{
		ClientId: creds.ClientID,
		UserId:   creds.UserID,
	})
	if err != nil {
		return nil, err
	}
	userDto := &response.UserResponse{
		TwitchId:        userInfo.Id,
		DisplayName:     userInfo.DisplayName,
		NameColor:       userColor.Color,
		ProfileImageUrl: userInfo.ProfileImageUrl,
	}
	return userDto, nil
}
