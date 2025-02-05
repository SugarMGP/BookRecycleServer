package studentController

import (
	"errors"

	"bookcycle-server/internal/apiException"
	"bookcycle-server/internal/services/studentService"
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

type loginResp struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	HasLogged bool   `json:"has_logged"`
}

// Login 学生登录
func Login(c *gin.Context) {
	var data loginReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
	}

	user, err := studentService.GetStudentByUsername(data.Username)
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

	zap.L().Info("学生登录成功", zap.String("username", data.Username))
	response.JsonSuccessResp(c, loginResp{
		Token:     token,
		Name:      user.Name,
		HasLogged: user.HasLogged,
	})
}
