package routes

import (
	"bookrecycle-server/internal/controllers/objectController"
	"bookrecycle-server/internal/controllers/userController"
	"bookrecycle-server/internal/midwares"
	"bookrecycle-server/pkg/ws"
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(r *gin.Engine) {
	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/login", userController.Login)
			user.POST("/register", userController.Register)
			user.POST("/activate", midwares.Auth, userController.Activate)
			user.GET("/info", midwares.Auth, userController.GetUserInfo)
		}
		api.POST("/upload", midwares.Auth, objectController.UploadFile)
	}
	r.GET("/ws", midwares.Auth, ws.HandleWebSocket)
}
