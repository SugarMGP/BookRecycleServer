package recycleController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/recycleService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// SettleRecycle 结算回收订单
func SettleRecycle(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	recycle, err := recycleService.GetRecycleByID(user.CurrentOrder)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	seller, err := userService.GetUserByID(recycle.Seller)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	seller.CurrentOrder = 0
	err = userService.SaveUser(seller)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	user.CurrentOrder = 0
	err = userService.SaveUser(user)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	recycle.Status = 4
	err = recycleService.SaveRecycle(recycle)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// TODO: 生成订单

	response.JsonSuccessResp(c, nil)
}
