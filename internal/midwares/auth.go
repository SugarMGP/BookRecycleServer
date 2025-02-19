package midwares

import (
	"errors"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/userService"
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
		if errors.Is(err, jwt.ErrTokenHandlingFailed) {
			response.AbortWithException(c, apiException.ServerError, err)
			return
		}
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
