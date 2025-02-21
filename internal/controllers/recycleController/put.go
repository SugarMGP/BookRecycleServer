package recycleController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/recycleService"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type putRecycleInfoReq struct {
	Img    string  `json:"img" binding:"required"`
	Weight float64 `json:"weight" binding:"required"`
}

// PutRecycleInfo 完善上门信息
func PutRecycleInfo(c *gin.Context) {
	var data putRecycleInfoReq
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

	recycle, err := recycleService.GetRecycleByID(user.CurrentOrder)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	recycle.Img = data.Img
	recycle.Weight = data.Weight
	recycle.Status = 3
	err = recycleService.SaveRecycle(recycle)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, nil)
}
