package recycleService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/pkg/database"
)

// SaveRecycle 保存回收订单
func SaveRecycle(recycle *models.Recycle) error {
	result := database.DB.Save(recycle)
	return result.Error
}
