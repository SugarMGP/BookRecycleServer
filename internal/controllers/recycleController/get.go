package recycleController

import (
	"time"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/recycleService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type getRecycleStatusResp struct {
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

	resp := getRecycleStatusResp{}

	if user.CurrentOrder == 0 {
		resp.Status = 0
	} else {
		order, err := recycleService.GetRecycleByID(user.CurrentOrder)
		if err != nil {
			response.AbortWithException(c, apiException.ServerError, err)
			return
		}
		resp.Status = order.Status

		if order.Status == 2 {
			receiver, err := userService.GetUserByID(order.Receiver)
			if err != nil {
				response.AbortWithException(c, apiException.ServerError, err)
				return
			}
			resp.ReceiverName = receiver.Name
			resp.ReceiverPhone = receiver.Phone
		}

		if order.Status == 3 {
			resp.Weight = order.Weight
			resp.Money = decimal.NewFromFloat(order.Weight * 1.6 * 0.8).StringFixedBank(2)
		}
	}

	response.JsonSuccessResp(c, resp)
}

type getCurrentOrderResp struct {
	ID          uint      `json:"id"`
	SellerName  string    `json:"seller_name"`
	SellerPhone string    `json:"seller_phone"`
	Weight      float64   `json:"weight"`
	Address     string    `json:"address"`
	Note        string    `json:"note"`
	Img         string    `json:"img"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetCurrentOrder 收书员获取当前订单
func GetCurrentOrder(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	resp := getCurrentOrderResp{}

	resp.ID = user.CurrentOrder
	if user.CurrentOrder == 0 {
		response.JsonSuccessResp(c, resp)
		return
	}

	order, err := recycleService.GetRecycleByID(user.CurrentOrder)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	resp.Img = order.Img
	resp.Address = order.Address
	resp.CreatedAt = order.CreatedAt
	resp.Note = order.Note
	resp.Weight = order.Weight
	resp.SellerName = order.SellerName

	seller, err := userService.GetUserByID(order.Seller)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}
	resp.SellerPhone = seller.Phone

	response.JsonSuccessResp(c, resp)
}

type orderListElement struct {
	ID         uint      `json:"id"`
	SellerName string    `json:"seller_name"`
	Weight     float64   `json:"weight"`
	Address    string    `json:"address"`
	Note       string    `json:"note"`
	Img        string    `json:"img"`
	CreatedAt  time.Time `json:"created_at"`
}

// GetOrderList 收书员获取列表
func GetOrderList(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	orders, err := recycleService.GetListByCampus(user.Campus)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	resp := make([]orderListElement, 0)
	for _, order := range orders {
		resp = append(resp, orderListElement{
			ID:         order.ID,
			SellerName: order.SellerName,
			Weight:     order.Weight,
			Address:    order.Address,
			Note:       order.Note,
			Img:        order.Img,
			CreatedAt:  order.CreatedAt,
		})
	}

	response.JsonSuccessResp(c, gin.H{
		"order_list": resp,
	})
}
