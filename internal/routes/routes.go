package routes

import (
	"bookcycle-server/internal/controllers/userController"
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(r *gin.Engine) {
	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", userController.Register)
			user.POST("/login", userController.Login)
		}
	}
}
