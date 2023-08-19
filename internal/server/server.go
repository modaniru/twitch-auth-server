package server

import (
	"github.com/gin-gonic/gin"
	"github.com/modaniru/twitch-auth-server/internal/dto/request"
	"github.com/modaniru/twitch-auth-server/internal/service"
)

// TODO error func
type MyServer struct {
	engine      *gin.Engine
	userService service.UserServicer
	authService service.AuthServicer
}

func NewMyServer(engine *gin.Engine, userService service.UserServicer, authService service.AuthServicer) *MyServer {
	s := &MyServer{
		engine:      engine,
		userService: userService,
		authService: authService,
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

func (m *MyServer) signIn(c *gin.Context) {
	token := new(request.AccessToken)
	err := c.ShouldBindJSON(token)
	if err != nil {
		c.JSON(403, err.Error())
		return
	}
	jwtToken, err := m.authService.Auth(token.Token)
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
	user, err := m.userService.GetUserInformation(id)
	if err != nil {
		c.JSON(403, err.Error())
		return
	}
	c.JSON(200, user)
}
