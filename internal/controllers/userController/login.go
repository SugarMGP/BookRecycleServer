package userController

import (
	"errors"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils/jwt"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Type     uint   `json:"type" binding:"required"`
}

type loginResp struct {
	Token     string `json:"token"`
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	Activated bool   `json:"activated"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var data loginReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	user, err := userService.GetUserByUsernameAndType(data.Username, data.Type)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.WrongPasswordOrUsername, err)
		return
	}
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 校验密码
	if user.Password != data.Password {
		response.AbortWithException(c, apiException.WrongPasswordOrUsername, err)
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, loginResp{
		Token:     token,
		UserID:    user.ID,
		Name:      user.Name,
		Activated: user.StudentID != "",
	})
}
