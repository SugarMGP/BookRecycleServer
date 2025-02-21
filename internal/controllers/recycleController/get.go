package recycleController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/recycleService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type getRecycleStatusReq struct {
	Status        uint    `json:"status"`
	ReceiverName  string  `json:"receiver_name"`
	ReceiverPhone string  `json:"receiver_phone"`
	Weight        float64 `json:"weight"`
	Money         string  `json:"money"`
}

// GetRecycleStatus 学生获取回收状态
func GetRecycleStatus(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	req := getRecycleStatusReq{}

	if user.CurrentOrder == 0 {
		req.Status = 0
	} else {
		order, err := recycleService.GetRecycleByID(user.CurrentOrder)
		if err != nil {
			response.AbortWithException(c, apiException.ServerError, err)
			return
		}
		req.Status = order.Status

		if order.Status == 2 {
			receiver, err := userService.GetUserByID(order.Receiver)
			if err != nil {
				response.AbortWithException(c, apiException.ServerError, err)
				return
			}
			req.ReceiverName = receiver.Name
			req.ReceiverPhone = receiver.Phone
		}

		if order.Status == 3 {
			req.Weight = order.Weight
			req.Money = decimal.NewFromFloat(order.Weight * 1.6 * 0.8).StringFixedBank(2)
		}
	}

	response.JsonSuccessResp(c, req)
}
