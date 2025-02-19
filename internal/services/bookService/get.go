package bookService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetBookList 获取书籍列表
func GetBookList(page, size int) ([]models.Book, int64, error) {
	var books []models.Book
	var totalRecords int64

	// 先查询总记录数
	if err := database.DB.Model(&models.Book{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// 计算总页数（向上取整）
	totalPages := (totalRecords + int64(size) - 1) / int64(size)

	// 执行分页查询
	result := database.DB.Offset((page - 1) * size).Limit(size).Find(&books)
	return books, totalPages, result.Error
}
