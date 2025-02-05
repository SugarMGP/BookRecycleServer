package userService

import (
	"bookcycle-server/internal/models"
	"bookcycle-server/pkg/database"
)

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}
