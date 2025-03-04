package reportService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetReportList 获取举报列表
func GetReportList() ([]models.Report, error) {
	var reports []models.Report
	result := database.DB.
		Order("CASE WHEN status = 1 THEN 0 ELSE 1 END").
		Order("id desc").
		Find(&reports)
	return reports, result.Error
}
