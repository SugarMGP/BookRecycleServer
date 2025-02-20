package feedbackController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/feedbackService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type feedbackReq struct {
	Content   string `json:"content" binding:"required"`
	Anonymity bool   `json:"anonymity"`
}

// Feedback 提交反馈
func Feedback(c *gin.Context) {
	var data feedbackReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	uid := user.ID
	if data.Anonymity {
		uid = 0
	}

	err = feedbackService.CreateFeedback(&models.Feedback{
		Content: data.Content,
		UserID:  uid,
	})
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
