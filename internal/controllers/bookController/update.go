package bookController

import (
	"errors"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/bookService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type updateBookReq struct {
	ID           uint   `json:"id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Author       string `json:"author" binding:"required"`
	Course       string `json:"course"`
	Edition      string `json:"edition"`
	Publisher    string `json:"publisher" binding:"required"`
	Completeness string `json:"completeness" binding:"required"`
	Img          string `json:"img" binding:"required"`
	Price        string `json:"price" binding:"required"`
	Note         string `json:"note"`
}

// UpdateBook 更新书籍
func UpdateBook(c *gin.Context) {
	var data updateBookReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	// 价格校验
	price, err := decimal.NewFromString(data.Price)
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

	{ // 更新书籍信息
		book.Name = data.Name
		book.Author = data.Author
		book.Course = data.Course
		book.Edition = data.Edition
		book.Publisher = data.Publisher
		book.Completeness = data.Completeness
		book.Img = data.Img
		book.Price = price.StringFixedBank(2)
		book.Note = data.Note
	}

	err = bookService.SaveBook(book)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
