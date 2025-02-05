package userController

import (
	"errors"

	"bookcycle-server/internal/apiException"
	"bookcycle-server/internal/services/userService"
	"bookcycle-server/internal/utils/jwt"
	"bookcycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var data loginReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
	}

	user, err := userService.GetUserByUsername(data.Username)
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

	zap.L().Info("用户登陆成功", zap.String("username", data.Username))
	response.JsonSuccessResp(c, gin.H{
		"token": token,
	})
}
