package recycleController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/billService"
	"bookrecycle-server/internal/services/recycleService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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

	money := recycle.Weight * 1.6
	student := decimal.NewFromFloat(money * 0.8)
	receiver := decimal.NewFromFloat(money * 0.2)
	// 保存学生账单
	err = billService.SaveBill(&models.Bill{
		User:   seller.ID,
		Amount: student.StringFixedBank(2),
	})
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 保存收书员账单
	err = billService.SaveBill(&models.Bill{
		User:   user.ID,
		Amount: receiver.StringFixedBank(2),
	})
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	seller.CurrentOrder = 0
	sellerBalance, _ := decimal.NewFromString(seller.Balance)
	seller.Balance = sellerBalance.Add(student).StringFixedBank(2)
	err = userService.SaveUser(seller)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	user.CurrentOrder = 0
	userBalance, _ := decimal.NewFromString(user.Balance)
	user.Balance = userBalance.Add(receiver).StringFixedBank(2)
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

	response.JsonSuccessResp(c, nil)
}
