package feedbackController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/feedbackService"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// GetFeedbackList 获取所有反馈
func GetFeedbackList(c *gin.Context) {
	// 获取所有反馈
	feedbackList, err := feedbackService.GetFeedbackList()
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, gin.H{
		"feedback_list": feedbackList,
	})
}
