package feedbackService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// SaveFeedback 保存反馈
func SaveFeedback(feedback *models.Feedback) error {
	result := database.DB.Save(feedback)
	return result.Error
}
