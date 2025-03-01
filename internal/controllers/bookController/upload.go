package bookController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/bookService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type uploadBookReq struct {
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

// UploadBook 上传书籍
func UploadBook(c *gin.Context) {
	var data uploadBookReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	price, err := decimal.NewFromString(data.Price)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	err = bookService.SaveBook(&models.Book{
		UserID:       user.ID,
		Name:         data.Name,
		Author:       data.Author,
		Course:       data.Course,
		Edition:      data.Edition,
		Publisher:    data.Publisher,
		Completeness: data.Completeness,
		Img:          data.Img,
		Price:        price.StringFixedBank(2),
		Note:         data.Note,
		Status:       3,
	})
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
