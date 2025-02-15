package userController

import (
	"errors"
	"regexp"

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
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Campus   uint   `json:"campus" binding:"required"`
	Type     uint   `json:"type" binding:"required"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var data registerReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	// 用户名校验
	if !isUsernameValid(data.Username) {
		response.AbortWithException(c, apiException.InvalidUsername, nil)
		return
	}

	// 密码校验
	if !isPasswordValid(data.Password) {
		response.AbortWithException(c, apiException.InvalidPassword, nil)
		return
	}

	// 校验校区信息
	if data.Campus > 3 || data.Campus < 1 {
		response.AbortWithException(c, apiException.ParamsError, nil)
		return
	}

	// 判断用户是否已经注册
	_, err = userService.GetUserByUsername(data.Username)
	if err == nil {
		response.AbortWithException(c, apiException.UserAlreadyExist, err)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 注册用户
	err = userService.SaveUser(models.User{
		Username: data.Username,
		Password: data.Password,
		Type:     data.Type,
		Name:     data.Name,
		Phone:    data.Phone,
		Campus:   data.Campus,
	})
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	zap.L().Info("用户注册成功",
		zap.String("username", data.Username),
		zap.Uint("type", data.Type),
	)
	response.JsonSuccessResp(c, nil)
}

// 用户名正则，4到16位（字母，数字，下划线，减号）
func isUsernameValid(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]{4,16}$`).MatchString(username)
}

// 密码正则，8到32位，至少一个字母和一个数字，可包含大小写字母、数字和特殊符号
func isPasswordValid(password string) bool {
	return regexp.MustCompile(`[A-Za-z]`).MatchString(password) &&
		regexp.MustCompile(`\d`).MatchString(password) &&
		regexp.MustCompile(`^[A-Za-z0-9\W]{8,32}$`).MatchString(password)
}
