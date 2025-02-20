package userController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type activateReq struct {
	Address   string `json:"address"`
	Campus    uint   `json:"campus"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	StudentID string `json:"student_id"`
}

// Activate 学生激活
func Activate(c *gin.Context) {
	var data activateReq
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

	// 判断校区信息是否合法
	if data.Campus > 3 || data.Campus < 1 {
		response.AbortWithException(c, apiException.ParamsError, nil)
		return
	}

	{ // 更新用户信息
		user.Address = data.Address
		user.Campus = data.Campus
		user.Name = data.Name
		user.Phone = data.Phone
		user.StudentID = data.StudentID
	}

	err = userService.SaveUser(user)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
