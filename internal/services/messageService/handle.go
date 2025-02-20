package messageService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetMessagesByUser 根据用户获取聊天消息
func GetMessagesByUser(user uint) ([]models.Message, error) {
	var messages []models.Message
	result := database.DB.Where("sender = ?", user).Or("receiver = ?", user).Find(&messages)
	return messages, result.Error
}

// CreateMessage 保存聊天消息
func CreateMessage(message *models.Message) error {
	result := database.DB.Create(message)
	return result.Error
}
