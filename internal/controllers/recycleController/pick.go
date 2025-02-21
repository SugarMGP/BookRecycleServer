package recycleController

import (
	"errors"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/recycleService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type pickRecycleReq struct {
	ID uint `json:"id" binding:"required"`
}

// PickRecycle 收书员接取回收订单
func PickRecycle(c *gin.Context) {
	var data pickRecycleReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	order, err := recycleService.GetRecycleByID(data.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.AbortWithException(c, apiException.ResourceNotFound, err)
		return
	}
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	order.Receiver = user.ID
	order.Status = 2
	err = recycleService.SaveRecycle(order)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	user.CurrentOrder = order.ID
	err = userService.SaveUser(user)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
