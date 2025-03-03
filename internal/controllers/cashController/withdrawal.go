package cashController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/cashService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type withdrawalReq struct {
	Amount float64 `json:"amount"`
}

func Withdrawal(c *gin.Context) {
	var req withdrawalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	if err = cashService.Withdrawal(user.ID, req.Amount); err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	response.JsonSuccessResp(c, nil)
}

func GetWithdrawalList(c *gin.Context) {
	withdrawalList, err := cashService.GetWithdrawalList()
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	response.JsonSuccessResp(c, withdrawalList)
}
