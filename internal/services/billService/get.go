package billService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetBillListByUser 根据用户ID获取账单列表
func GetBillListByUser(uid uint) ([]models.Bill, error) {
	var bills []models.Bill
	result := database.DB.Where("user = ?", uid).Find(&bills)
	return bills, result.Error
}
