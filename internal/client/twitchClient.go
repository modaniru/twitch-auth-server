package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	validate         = "https://id.twitch.tv/oauth2/validate"
	getUsers         = "https://api.twitch.tv/helix/users"
	getUserChatColor = "https://api.twitch.tv/helix/chat/color"
	getAppToken      = "https://id.twitch.tv/oauth2/token"
)

type TwitchClienter interface {
	Validate(token string) (*ValidateResponse, error)
	GetUserInformation(token string, userCred *ValidateResponse) (*UserInformation, error)
	GetUserInformationApp(userCred *ValidateResponse) (*UserInformation, error)
	GetUserChatColor(token string, userCred *ValidateResponse) (*UserColor, error)
	GetUserChatColorApp(userCred *ValidateResponse) (*UserColor, error)
}

type TwitchClient struct {
	client             *Client
	twitchClientId     string
	twitchClientSecret string
	token              string
}

func NewTwitchClient(client *Client, twitchClientId, twitchClientSecret string) TwitchClienter {
	return &TwitchClient{client: client, twitchClientId: twitchClientId, twitchClientSecret: twitchClientSecret}
}

// TODO dry, create my client
func (t *TwitchClient) Validate(token string) (*ValidateResponse, error) {
	headers := map[string]string{
		"Authorization": "OAuth " + token,
	}
	response, err := t.client.Request("GET", validate, headers, nil, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		bytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("status code: %d, uri %s. Error: %s", response.StatusCode, validate, string(bytes))
	}
	body := response.Body
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	validate := new(ValidateResponse)
	err = json.Unmarshal(bytes, validate)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

func (t *TwitchClient) GetUserInformationApp(userCred *ValidateResponse) (*UserInformation, error) {
	err := t.getToken()
	if err != nil {
		return nil, err
	}
	return t.GetUserInformation(t.token, userCred)
}

func (t *TwitchClient) GetUserChatColorApp(userCred *ValidateResponse) (*UserColor, error) {
	err := t.getToken()
	if err != nil {
		return nil, err
	}
	return t.GetUserChatColor(t.token, userCred)
}

func (t *TwitchClient) GetUserInformation(token string, userCred *ValidateResponse) (*UserInformation, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Client-Id":     t.twitchClientId,
	}
	params := map[string]string{
		"id": userCred.UserId,
	}
	response, err := t.client.Request("GET", getUsers, headers, params, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		bytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("status code: %d, uri %s. Error: %s", response.StatusCode, getUsers, string(bytes))
	}
	body := response.Body
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	usersInfo := new(UsersInformation)
	err = json.Unmarshal(bytes, usersInfo)
	if err != nil {
		return nil, err
	}
	if len(usersInfo.Data) == 0 {
		return nil, errors.New("len data UsersInformation equal zero")
	}
	return &usersInfo.Data[0], nil
}

func (t *TwitchClient) GetUserChatColor(token string, userCred *ValidateResponse) (*UserColor, error) {
	t.getToken()
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Client-Id":     t.twitchClientId,
	}
	params := map[string]string{
		"user_id": userCred.UserId,
	}
	response, err := t.client.Request("GET", getUserChatColor, headers, params, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		bytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("status code: %d, uri %s. Error: %s", response.StatusCode, getUserChatColor, string(bytes))
	}
	body := response.Body
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	usersColor := new(UsersColor)
	err = json.Unmarshal(bytes, usersColor)
	if err != nil {
		return nil, err
	}
	if len(usersColor.Data) == 0 {
		return nil, errors.New("len data UsersColor equal zero")
	}
	return usersColor.Data[0], nil
}

func (t *TwitchClient) getToken() error {
	headers := map[string]string{
		"Authorization": "OAuth " + t.token,
	}
	response, err := t.client.Request("GET", validate, headers, nil, nil)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		params := map[string]string{
			"client_id":     t.twitchClientId,
			"client_secret": t.twitchClientSecret,
			"grant_type":    "client_credentials",
		}
		response, err := t.client.Request("POST", getAppToken, nil, params, nil)
		if err != nil {
			return nil
		}
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil
		}
		object := new(AppCredentials)
		err = json.Unmarshal(bytes, object)
		if err != nil {
			return nil
		}
		t.token = object.AccessToken
	}
	return nil
}
