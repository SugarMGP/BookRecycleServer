package bookController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/bookService"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type getBookListReq struct {
	Page int `json:"page" binding:"required"`
	Size int `json:"size" binding:"required"`
}

type getBookListResp struct {
	BookList   []bookListElement `json:"book_list"`
	TotalPages int64             `json:"total_pages"`
}

type bookListElement struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	Course       string `json:"course"`
	Edition      string `json:"edition"`
	Publisher    string `json:"publisher"`
	Completeness string `json:"completeness"`
	Img          string `json:"img"`
	Price        string `json:"price"`
	Note         string `json:"note"`
}

// GetBookList 获取书籍列表
func GetBookList(c *gin.Context) {
	var data getBookListReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	books, totalPages, err := bookService.GetBookList(data.Page, data.Size)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	bookList := make([]bookListElement, 0)
	for _, book := range books {
		bookList = append(bookList, bookListElement{
			ID:           book.ID,
			UserID:       book.UserID,
			Name:         book.Name,
			Course:       book.Course,
			Edition:      book.Edition,
			Publisher:    book.Publisher,
			Completeness: book.Completeness,
			Img:          book.Img,
			Price:        book.Price,
			Note:         book.Note,
		})
	}

	response.JsonSuccessResp(c, getBookListResp{
		BookList:   bookList,
		TotalPages: totalPages,
	})
}
