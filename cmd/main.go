package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/modaniru/twitch-auth-server/internal/client"
	"github.com/modaniru/twitch-auth-server/internal/db"
	"github.com/modaniru/twitch-auth-server/internal/repository"
	"github.com/modaniru/twitch-auth-server/internal/server"
	"github.com/modaniru/twitch-auth-server/internal/service"
	_ "github.com/lib/pq"
)

//TODO migrations
//TODO remove sqlc
func main() {
	_ = godotenv.Load()
	dns := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	database, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err.Error())
	}
	userRepository := repository.NewUserRepository(database, db.New(database))
	twitchClient := client.NewTwitchClient(client.NewClient(http.Client{}), os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	userService := service.NewUserService(userRepository, twitchClient)
	authService := service.NewAuthService(userService)
	server := server.NewMyServer(gin.Default(), userService, authService)
	server.Run("8080")
}
