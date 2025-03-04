package reportService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// SaveReport 保存举报
func SaveReport(report *models.Report) error {
	result := database.DB.Save(report)
	return result.Error
}
