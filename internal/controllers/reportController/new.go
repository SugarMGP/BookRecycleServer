package reportController

import (
	"errors"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/bookService"
	"bookrecycle-server/internal/services/reportService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type reportReq struct {
	BookID uint   `json:"book_id" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

// NewReport 新建举报
func NewReport(c *gin.Context) {
	var data reportReq
	if err := c.ShouldBindJSON(&data); err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	book, err := bookService.GetBookByID(data.BookID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.ResourceNotFound, err)
		return
	}
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	seller, err := userService.GetUserByID(book.UserID)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	err = reportService.SaveReport(&models.Report{
		Reporter:     user.ID,
		ReporterName: user.Name,
		Seller:       seller.ID,
		SellerName:   seller.Name,
		Book:         data.BookID,
		BookName:     book.Name,
		Title:        data.Title,
		Status:       1,
	})
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
