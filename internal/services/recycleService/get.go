package recycleService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// GetRecycleByID 通过ID获取回收订单
func GetRecycleByID(id uint) (*models.Recycle, error) {
	var recycle models.Recycle
	result := database.DB.Where("id = ?", id).First(&recycle)
	return &recycle, result.Error
}
