package bookService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// DeleteBook 删除书籍
func DeleteBook(id uint) error {
	result := database.DB.Where("id = ?", id).Delete(&models.Book{})
	return result.Error
}
