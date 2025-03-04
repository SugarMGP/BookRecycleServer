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

// GetReportByID 根据ID获取举报
func GetReportByID(id uint) (*models.Report, error) {
	var report models.Report
	result := database.DB.Where("id = ?", id).First(&report)
	return &report, result.Error
}
