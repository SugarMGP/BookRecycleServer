package studentService

import (
	"bookcycle-server/internal/models"
	"bookcycle-server/pkg/database"
)

// GetStudentByUsername 通过学号获取学生
func GetStudentByUsername(username string) (*models.Student, error) {
	var student models.Student
	result := database.DB.Where("username = ?", username).First(&student)
	return &student, result.Error
}
