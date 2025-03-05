package captchaController

import (
	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/utils/response"
	"bookrecycle-server/pkg/captcha"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
)

type getCaptchaResp struct {
	CaptchaKey  string `json:"captcha_key"`
	ImageBase64 string `json:"image_base64"`
	TileBase64  string `json:"tile_base64"`
	TileWidth   int    `json:"tile_width"`
	TileHeight  int    `json:"tile_height"`
	TileX       int    `json:"tile_x"`
	TileY       int    `json:"tile_y"`
}

// GetCaptcha 获取验证码
func GetCaptcha(c *gin.Context) {
	captData, err := captcha.SlideCaptcha.Generate()
	if err != nil {
		response.AbortWithException(c, apiException.CaptchaError, err)
		return
	}

	blockData := captData.GetData()
	uid := uuid.NewV1().String()
	err = captcha.CaptchaStore.Add(uid, blockData, cache.DefaultExpiration)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		response.AbortWithException(c, apiException.CaptchaError, err)
		return
	}
	tBase64, err = captData.GetTileImage().ToBase64()
	if err != nil {
		response.AbortWithException(c, apiException.CaptchaError, err)
		return
	}

	response.JsonSuccessResp(c, getCaptchaResp{
		CaptchaKey:  uid,
		ImageBase64: mBase64,
		TileBase64:  tBase64,
		TileWidth:   blockData.Width,
		TileHeight:  blockData.Height,
		TileX:       blockData.TileX,
		TileY:       blockData.TileY,
	})
}
