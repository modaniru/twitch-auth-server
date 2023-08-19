package server

import "github.com/gin-gonic/gin"

func (m *MyServer) auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	id, err := m.authService.Validate(token)
	if err != nil {
		c.JSON(403, err.Error())
		c.Abort()
		return
	}
	c.Set("id", id)
}
