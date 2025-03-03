package midwares

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/jwt"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// Auth 用户认证中间件
func Auth(usertype ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过 header 中的 Authorization 来认证
		token := c.Request.Header.Get("Authorization")
		if token == "" { // 没有携带 token
			response.AbortWithException(c, apiException.NoAccessPermission, nil)
			return
		}

		token = token[7:] // 去除 Bearer
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

		if !user.Activated && user.Type != 3 {
			response.AbortWithException(c, apiException.UserNotActive, nil)
			return
		}

		// 判断用户类型
		if len(usertype) != 0 {
			flag := false
			for _, ut := range usertype {
				if ut == user.Type {
					flag = true
					break
				}
			}
			if !flag {
				response.AbortWithException(c, apiException.NoAccessPermission, nil)
				return
			}
		}

		c.Set("user", user)
		c.Next()
	}
}

// AuthReviewBooks 书籍审核权限检查
func AuthReviewBooks(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	if !user.CanReviewBooks {
		response.AbortWithException(c, apiException.NoAccessPermission, nil)
		return
	}

	c.Next()
}

// AuthManageReports 举报管理权限检查
func AuthManageReports(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	if !user.CanManageReports {
		response.AbortWithException(c, apiException.NoAccessPermission, nil)
		return
	}

	c.Next()
}
