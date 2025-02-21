package recycleController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/recycleService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type uploadRecycleReq struct {
	Img     string  `json:"img" binding:"required"`
	Note    string  `json:"note"`
	Address string  `json:"address" binding:"required"`
	Weight  float64 `json:"weight" binding:"required"`
}

// UploadRecycle 学生提交回收请求
func UploadRecycle(c *gin.Context) {
	var data uploadRecycleReq
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

	recycle := models.Recycle{
		Seller:     user.ID,
		SellerName: user.Name,
		Img:        data.Img,
		Note:       data.Note,
		Weight:     data.Weight,
		Address:    data.Address,
		Campus:     user.Campus,
		Status:     1,
	}
	err = recycleService.SaveRecycle(&recycle)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	user.CurrentOrder = recycle.ID
	err = userService.SaveUser(user)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
