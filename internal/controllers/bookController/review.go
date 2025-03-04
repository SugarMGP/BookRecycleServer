package bookController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/bookService"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type getReviewBookListReq struct {
	Search string `json:"search"`
	Status uint   `json:"status"`
}

type reviewBookListElement struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	Author       string `json:"author"`
	Course       string `json:"course"`
	Edition      string `json:"edition"`
	Publisher    string `json:"publisher"`
	Completeness string `json:"completeness"`
	Img          string `json:"img"`
	Price        string `json:"price"`
	Note         string `json:"note"`
	Status       uint   `json:"status"`
	Reason       string `json:"reason"`
}

// GetReviewBookList 获取书籍审核列表
func GetReviewBookList(c *gin.Context) {
	var data getReviewBookListReq
	err := c.ShouldBind(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	books, err := bookService.GetReviewBookList(data.Search, data.Status)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	bookList := make([]reviewBookListElement, 0)
	for _, book := range books {
		bookList = append(bookList, reviewBookListElement{
			ID:           book.ID,
			UserID:       book.UserID,
			Name:         book.Name,
			Author:       book.Author,
			Course:       book.Course,
			Edition:      book.Edition,
			Publisher:    book.Publisher,
			Completeness: book.Completeness,
			Img:          book.Img,
			Price:        book.Price,
			Note:         book.Note,
			Status:       book.Status,
			Reason:       book.Reason,
		})
	}

	response.JsonSuccessResp(c, gin.H{
		"review_book_list": bookList,
	})
}

type updateReviewStatusReq struct {
	ID     uint   `json:"id" binding:"required"`
	Status uint   `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

// UpdateReviewStatus 更新书籍审核状态
func UpdateReviewStatus(c *gin.Context) {
	var data updateReviewStatusReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	if data.Status != 1 && data.Status != 4 {
		response.AbortWithException(c, apiException.ParamsError, nil)
		return
	}

	book, err := bookService.GetBookByID(data.ID)
	book.Status = data.Status
	book.Reason = data.Reason

	err = bookService.SaveBook(book)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
