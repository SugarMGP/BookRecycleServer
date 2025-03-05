package userController

import (
	"time"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/billService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type billElement struct {
	Time  string `json:"time"`
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
	Bills      []billElement `json:"bills"`
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	bills := make([]billElement, 0)
	billList, err := billService.GetBillListByUser(user.ID)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	for _, bill := range billList {
		bills = append(bills, billElement{
			Time:  bill.CreatedAt.Format(time.DateTime),
			Money: bill.Amount,
		})
	}

	response.JsonSuccessResp(c, infoResp{
		ID:         user.ID,
		Name:       user.Name,
		StudentID:  user.StudentID,
		Phone:      user.Phone,
		Campus:     user.Campus,
		Address:    user.Address,
		Balance:    user.Balance,
		Reputation: user.Reputation,
		Bills:      bills,
	})
}
