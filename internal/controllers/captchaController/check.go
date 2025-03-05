package captchaController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/utils/response"
	"bookrecycle-server/pkg/captcha"
	"github.com/gin-gonic/gin"
	"github.com/wenlng/go-captcha/v2/slide"
)

type checkCaptchaReq struct {
	CaptchaKey string `json:"captcha_key"  binding:"required"`
	X          int64  `json:"x"  binding:"required"`
	Y          int64  `json:"y"  binding:"required"`
}

// CheckCaptcha 校验验证码
func CheckCaptcha(c *gin.Context) {
	var data checkCaptchaReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	block, ok1 := captcha.CaptchaStore.Get(data.CaptchaKey)
	dct, ok2 := block.(*slide.Block)
	if !ok1 || !ok2 {
		response.AbortWithException(c, apiException.CaptchaTimeout, nil)
		return
	}

	flag := slide.CheckPoint(data.X, data.Y, int64(dct.X), int64(dct.Y), 4)
	response.JsonSuccessResp(c, gin.H{
		"check": flag,
	})

	captcha.CaptchaStore.Delete(data.CaptchaKey)
}
