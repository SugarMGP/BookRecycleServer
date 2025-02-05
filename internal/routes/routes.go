package routes

import (
	"bookcycle-server/internal/controllers/studentController"
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(r *gin.Engine) {
	api := r.Group("/api")
	{
		student := api.Group("/student")
		{
			student.POST("/login", studentController.Login)
		}
	}
}
