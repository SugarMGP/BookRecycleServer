package routes

import (
	"bookrecycle-server/internal/controllers/bookController"
	"bookrecycle-server/internal/controllers/captchaController"
	"bookrecycle-server/internal/controllers/cashController"
	"bookrecycle-server/internal/controllers/feedbackController"
	"bookrecycle-server/internal/controllers/objectController"
	"bookrecycle-server/internal/controllers/recycleController"
	"bookrecycle-server/internal/controllers/reportController"
	"bookrecycle-server/internal/controllers/userController"
	"bookrecycle-server/internal/midwares"
	"bookrecycle-server/pkg/ws"
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 用户接口
		user := api.Group("/user")
		{
			user.POST("/login", userController.Login)
			user.POST("/register", userController.Register)
			user.POST("/activate", userController.Activate)
			user.GET("/info", midwares.Auth(1, 2), userController.GetUserInfo)
			user.POST("/withdrawal", midwares.Auth(1, 2), cashController.Withdrawal)
		}

		// 学生接口
		student := api.Group("/student", midwares.Auth(1))
		{
			market := student.Group("/market")
			{
				market.POST("/products", bookController.GetBookList)
				market.GET("/books", bookController.GetMyBookList)
				market.POST("/book", bookController.UploadBook)
				market.PUT("/book", bookController.UpdateBook)
				market.DELETE("/book", bookController.DeleteBook)
				market.POST("/report", reportController.NewReport)
			}
			student.GET("/recycle", recycleController.GetRecycleStatus)
			student.POST("/recycle", recycleController.UploadRecycle)
			student.POST("/feedback", feedbackController.Feedback)
		}

		// 收书员接口
		receiver := api.Group("/receiver", midwares.Auth(2))
		{
			receiver.POST("/order", recycleController.PickRecycle)
			receiver.PUT("/order", recycleController.PutRecycleInfo)
			receiver.POST("/settle", recycleController.SettleRecycle)
			receiver.GET("/current_order", recycleController.GetCurrentOrder)
			receiver.GET("/orders", recycleController.GetOrderList)
		}

		// 管理员接口
		admin := api.Group("/admin", midwares.Auth(3))
		{
			admin.GET("/feedbacks", feedbackController.GetFeedbackList)
			admin.GET("/withdrawal", cashController.GetWithdrawalList)
			review := admin.Group("/review")
			{
				review.POST("/books", bookController.GetReviewBookList)
				review.PUT("/update", bookController.UpdateReviewStatus)
			}
			report := admin.Group("/report")
			{
				report.GET("/list", reportController.GetReportList)
				report.POST("/pass", reportController.PassReport)
				report.POST("/undo", reportController.UndoReport)
			}
		}

		api.POST("/upload", midwares.Auth(), objectController.UploadFile)

		// 验证码接口
		api.GET("/captcha", captchaController.GetCaptcha)
		api.POST("/captcha", captchaController.CheckCaptcha)
	}
	r.GET("/ws", ws.HandleWebSocket)
}
