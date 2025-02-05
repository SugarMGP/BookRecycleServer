package studentService

import (
	"bookcycle-server/internal/models"
	"bookcycle-server/pkg/database"
)

// SaveStudent 创建用户
func SaveStudent(student models.Student) error {
	result := database.DB.Save(&student)
	return result.Error
}
