package services

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/modaniru/twitch-auth-server/internal/service"
)

type AuthService struct {
	userService service.User
	secret      []byte
}

func NewAuthService(userService service.User, secret string) *AuthService {
	return &AuthService{
		userService: userService,
		secret:      []byte(secret),
	}
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
	return claims.SignedString(a.secret)
}

func (a *AuthService) Validate(jwtToken string) (int, error) {
	args := strings.Split(jwtToken, " ")
	jwtToken = args[1]
	t, err := jwt.ParseWithClaims(jwtToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(a.secret), nil
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
