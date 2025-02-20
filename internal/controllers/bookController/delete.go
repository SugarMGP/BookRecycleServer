package bookController

import (
	"errors"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/bookService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type deleteBookReq struct {
	ID uint `json:"id" binding:"required"`
}

// DeleteBook 下架书籍
func DeleteBook(c *gin.Context) {
	var data deleteBookReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	book, err := bookService.GetBookByID(data.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.ResourceNotFound, err)
		return
	}
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 判断书籍是否属于当前用户
	if book.UserID != user.ID {
		response.AbortWithException(c, apiException.NoAccessPermission, nil)
		return
	}

	err = bookService.DeleteBook(data.ID)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
