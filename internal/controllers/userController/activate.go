package userController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils/jwt"
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

	// 判断校区信息是否合法
	if data.Campus > 3 || data.Campus < 1 {
		response.AbortWithException(c, apiException.ParamsError, nil)
		return
	}

	token := c.Request.Header.Get("Authorization")
	token = token[7:]
	claims, err := jwt.ParseToken(token)
	if err != nil {
		response.AbortWithException(c, apiException.NoAccessPermission, err)
		return
	}

	// 获取用户信息
	user, err := userService.GetUserByID(claims.UserID)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 判断用户类型
	if user.Type != 1 && user.Type != 2 {
		response.AbortWithException(c, apiException.NoAccessPermission, nil)
		return
	}

	{ // 更新用户信息
		user.Address = data.Address
		user.Campus = data.Campus
		user.Name = data.Name
		user.Phone = data.Phone
		user.StudentID = data.StudentID
		user.Reputation = 100
		user.Balance = "0"
		user.Activated = true
	}

	err = userService.SaveUser(user)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
