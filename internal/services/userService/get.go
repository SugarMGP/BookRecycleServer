package userService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetUserByUsernameAndType 通过用户名和类型获取用户
func GetUserByUsernameAndType(username string, usertype uint) (*models.User, error) {
	var user models.User
	result := database.DB.Where("type = ?", usertype).Where("username = ?", username).First(&user)
	return &user, result.Error
}

// GetUserByID 通过ID获取用户
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", id).First(&user)
	return &user, result.Error
}
