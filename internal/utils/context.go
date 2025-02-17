package utils

import (
	"errors"

	"bookrecycle-server/internal/models"
	"github.com/gin-gonic/gin"
)

// GetUser 从上下文中提取 *models.User
func GetUser(c *gin.Context) (*models.User, error) {
	if val, ok := c.Get("user"); ok {
		if user, ok := val.(*models.User); ok {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
