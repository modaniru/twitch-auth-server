package services

import (
	"github.com/modaniru/twitch-auth-server/internal/client"
	"github.com/modaniru/twitch-auth-server/internal/entity"
	"github.com/modaniru/twitch-auth-server/internal/storage"
	"github.com/modaniru/twitch-auth-server/internal/storage/repo"
)

type UserService struct {
	userRepository storage.User
	twitchClient   client.TwitchClienter
}

// TODO construct
func NewUserService(userRepository storage.User, twitchClient client.TwitchClienter) *UserService {
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
	user, err := u.userRepository.GetUserByUserId(creds.UserId)
	if err == repo.ErrUserNotFound {
		return u.userRepository.CreateUser(creds.UserId, creds.ClientId)
	}
	return user.ID, nil
}

func (u *UserService) GetUserInformation(id int) (*entity.UserResponse, error) {
	creds, err := u.userRepository.GetUserById(id)
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
	userDto := &entity.UserResponse{
		TwitchId:        userInfo.Id,
		DisplayName:     userInfo.DisplayName,
		NameColor:       userColor.Color,
		ProfileImageUrl: userInfo.ProfileImageUrl,
	}
	return userDto, nil
}
