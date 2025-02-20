package feedbackService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetFeedbackList 获取反馈列表
func GetFeedbackList() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	result := database.DB.Order("id desc").Find(&feedbacks)
	return feedbacks, result.Error
}
