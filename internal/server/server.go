package server

import (
	"github.com/gin-gonic/gin"
	"github.com/modaniru/twitch-auth-server/internal/service"
)

// TODO error func
type MyServer struct {
	engine  *gin.Engine
	service *service.Service
}

func NewMyServer(engine *gin.Engine, service *service.Service) *MyServer {
	s := &MyServer{
		engine:  engine,
		service: service,
	}
	s.initRoutes()
	return s
}

func (m *MyServer) Run(port string) error {
	return m.engine.Run(":" + port)
}

func (m *MyServer) initRoutes() {
	m.engine.POST("/sign-in", m.signIn)
	api := m.engine.Group("/api", m.auth)
	{
		api.GET("/user", m.getUser)
	}
}

type accessToken struct {
	Token string `json:"token"`
}

func (m *MyServer) signIn(c *gin.Context) {
	token := new(accessToken)
	err := c.ShouldBindJSON(token)
	if err != nil {
		c.JSON(403, err.Error())
		return
	}
	jwtToken, err := m.service.AuthService.Auth(token.Token)
	if err != nil {
		c.JSON(403, err.Error())
		return
	}
	c.JSON(200, map[string]string{
		"jwt": jwtToken,
	})
}

func (m *MyServer) getUser(c *gin.Context) {
	id := c.GetInt("id")
	user, err := m.service.User.GetUserInformation(id)
	if err != nil {
		c.JSON(403, err.Error())
		return
	}
	c.JSON(200, user)
}
