package billService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// SaveBill 保存账单
func SaveBill(bill *models.Bill) error {
	result := database.DB.Save(bill)
	return result.Error
}
