package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/modaniru/twitch-auth-server/src/client"
	"github.com/modaniru/twitch-auth-server/src/db"
	"github.com/modaniru/twitch-auth-server/src/repository"
	"github.com/modaniru/twitch-auth-server/src/server"
	"github.com/modaniru/twitch-auth-server/src/service"

	_ "github.com/lib/pq"
)

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
