package feedbackService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// CreateFeedback 创建反馈
func CreateFeedback(feedback *models.Feedback) error {
	result := database.DB.Create(feedback)
	return result.Error
}
