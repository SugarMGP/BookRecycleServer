package cashController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/cashService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type withdrawalReq struct {
	Amount float64 `json:"amount" binding:"required"`
}

// Withdrawal 提现接口
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

	balance, _ := decimal.NewFromString(user.Balance)
	cash := decimal.NewFromFloat(req.Amount).Abs()
	newBalance := balance.Sub(cash)
	if newBalance.IsNegative() {
		response.AbortWithException(c, apiException.BalanceNotEnough, nil)
		return
	}

	if err = cashService.Withdrawal(user, newBalance, cash); err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	response.JsonSuccessResp(c, nil)
}

// GetWithdrawalList 获取提现记录
func GetWithdrawalList(c *gin.Context) {
	withdrawalList, err := cashService.GetWithdrawalList()
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	response.JsonSuccessResp(c, gin.H{
		"withdrawal_list": withdrawalList,
	})
}
