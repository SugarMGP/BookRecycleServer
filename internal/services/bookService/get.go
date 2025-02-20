package bookService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetBookList 获取书籍列表
func GetBookList(search string) ([]models.Book, error) {
	var books []models.Book
	query := database.DB.Where("status = 1")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ?", searchPattern).
			Or("course LIKE ?", searchPattern).
			Or("author LIKE ?", searchPattern)
	}

	result := query.Order("id desc").Find(&books)
	return books, result.Error
}

// GetMyBookList 获取我的书籍列表
func GetMyBookList(uid uint) ([]models.Book, error) {
	var books []models.Book
	result := database.DB.Where("user_id = ?", uid).Order("id desc").Find(&books)
	return books, result.Error
}

// GetBookByID 根据ID获取书籍
func GetBookByID(id uint) (*models.Book, error) {
	var book models.Book
	result := database.DB.Where("id = ?", id).First(&book)
	return &book, result.Error
}
