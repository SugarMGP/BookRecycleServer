package bookController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/bookService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

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
	books, err := bookService.GetBookList()
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

	response.JsonSuccessResp(c, gin.H{
		"book_list": bookList,
	})
}

type myBookListElement struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Course       string `json:"course"`
	Edition      string `json:"edition"`
	Publisher    string `json:"publisher"`
	Completeness string `json:"completeness"`
	Img          string `json:"img"`
	Price        string `json:"price"`
	Note         string `json:"note"`
	Status       uint   `json:"status"`
}

// GetMyBookList 获取我的书籍列表
func GetMyBookList(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	books, err := bookService.GetMyBookList(user.ID)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	bookList := make([]myBookListElement, 0)
	for _, book := range books {
		bookList = append(bookList, myBookListElement{
			ID:           book.ID,
			Name:         book.Name,
			Course:       book.Course,
			Edition:      book.Edition,
			Publisher:    book.Publisher,
			Completeness: book.Completeness,
			Img:          book.Img,
			Price:        book.Price,
			Note:         book.Note,
			Status:       book.Status,
		})
	}

	response.JsonSuccessResp(c, gin.H{
		"my_book_list": bookList,
	})
}
