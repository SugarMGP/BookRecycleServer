package userController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type billElement struct {
	Date  string `json:"date"`
	Money string `json:"money"`
}

type infoResp struct {
	ID         uint          `json:"id"`
	Name       string        `json:"name"`
	StudentID  string        `json:"student_id"`
	Phone      string        `json:"phone"`
	Campus     uint          `json:"campus"`
	Address    string        `json:"address"`
	Balance    string        `json:"balance"`
	Reputation uint          `json:"reputation"`
	Bill       []billElement `json:"bill"`
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// TODO 钱款相关信息
	response.JsonSuccessResp(c, infoResp{
		ID:         user.ID,
		Name:       user.Name,
		StudentID:  user.StudentID,
		Phone:      user.Phone,
		Campus:     user.Campus,
		Address:    user.Address,
		Balance:    "",
		Reputation: 0,
		Bill:       nil,
	})
}
