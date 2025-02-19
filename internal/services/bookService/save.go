package bookService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// SaveBook 保存书籍信息
func SaveBook(book models.Book) error {
	result := database.DB.Save(&book)
	return result.Error
}
