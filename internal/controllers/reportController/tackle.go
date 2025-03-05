package reportController

import (
	"errors"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/reportService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tackleReq struct {
	ID uint `json:"id" binding:"required"`
}

// PassReport 通过举报
func PassReport(c *gin.Context) {
	var req tackleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	report, err := reportService.GetReportByID(req.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.ResourceNotFound, err)
		return
	}
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 如果举报未通过，则降低用户信誉值
	if report.Status != 3 {
		user, err := userService.GetUserByID(report.Seller)
		if err != nil {
			response.AbortWithException(c, apiException.ServerError, err)
			return
		}

		user.Reputation -= 25
		err = userService.SaveUser(user)
		if err != nil {
			response.AbortWithException(c, apiException.ServerError, err)
			return
		}
	}

	report.Status = 3
	err = reportService.SaveReport(report)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}

// UndoReport 撤销举报
func UndoReport(c *gin.Context) {
	var req tackleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	report, err := reportService.GetReportByID(req.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.ResourceNotFound, err)
		return
	}
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 如果举报已通过，则撤销，并恢复用户信誉值
	if report.Status == 3 {
		user, err := userService.GetUserByID(report.Seller)
		if err != nil {
			response.AbortWithException(c, apiException.ServerError, err)
			return
		}

		user.Reputation += 25
		err = userService.SaveUser(user)
		if err != nil {
			response.AbortWithException(c, apiException.ServerError, err)
			return
		}
	}

	report.Status = 2
	err = reportService.SaveReport(report)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
