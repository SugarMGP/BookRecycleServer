package userController

import (
	"errors"

	"bookcycle-server/internal/apiException"
	"bookcycle-server/internal/models"
	"bookcycle-server/internal/services/userService"
	"bookcycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 注册
func Register(c *gin.Context) {
	var data registerReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
	}

	_, err = userService.GetUserByUsername(data.Username)
	if err == nil {
		response.AbortWithException(c, apiException.UserExisted, err)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	err = userService.SaveUser(models.User{
		Username: data.Username,
		Password: data.Password,
	})
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	zap.L().Info("用户创建成功", zap.String("username", data.Username))
	response.JsonSuccessResp(c, nil)
}
