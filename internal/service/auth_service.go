package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthServicer interface {
	// Return jwt token
	Auth(accessToken string) (string, error)
	Validate(jwt string) (int, error)
}

var (
	secret []byte = []byte("secret")
)

type AuthService struct {
	userService UserServicer
}

func NewAuthService(userService UserServicer) AuthServicer {
	return &AuthService{userService: userService}
}

func (a *AuthService) Auth(accessToken string) (string, error) {
	id, err := a.userService.CreateOrGetUser(accessToken)
	if err != nil {
		return "", err
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(12 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	return claims.SignedString(secret)
}

func (a *AuthService) Validate(jwtToken string) (int, error) {
	args := strings.Split(jwtToken, " ")
	jwtToken = args[1]
	t, err := jwt.ParseWithClaims(jwtToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("claims map error")
	}
	mapClaims := map[string]interface{}(claims)
	id := int(mapClaims["id"].(float64))
	return id, nil
}
