package userService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// SaveUser 保存用户
func SaveUser(user models.User) error {
	result := database.DB.Save(&user)
	return result.Error
}
