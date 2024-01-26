package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/modaniru/twitch-auth-server/internal/client"
	"github.com/modaniru/twitch-auth-server/internal/server"
	"github.com/modaniru/twitch-auth-server/internal/service"
	"github.com/modaniru/twitch-auth-server/internal/service/services"
	"github.com/modaniru/twitch-auth-server/internal/storage"
)

// TODO migrations
// TODO service structure in service package
// TODO handler to another package
// TODO slog
func main() {
	_ = godotenv.Load()
	dns := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	database, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err.Error())
	}
	userRepository := storage.NewStorage(database)
	twitchClient := client.NewTwitchClient(client.NewClient(http.Client{}), os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	userService := services.NewUserService(userRepository, twitchClient)
	authService := services.NewAuthService(userService, os.Getenv("JWT_SALT"))
	service := service.NewService(
		service.Dependencies{
			AuthService: authService,
			User:        userService,
		},
	)
	server := server.NewMyServer(gin.Default(), service)
	server.Run("8080")
}
