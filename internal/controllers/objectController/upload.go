package objectController

import (
	"errors"
	"image"
	"mime/multipart"
	"path"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/objectService"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type uploadFileData struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	var data uploadFileData
	if err := c.ShouldBind(&data); err != nil {
		response.AbortWithException(c, apiException.ParamsError, err)
		return
	}

	fileSize := data.File.Size
	if fileSize > objectService.SizeLimit {
		response.AbortWithException(c, apiException.FileSizeExceedError, nil)
		return
	}

	file, err := data.File.Open()
	if err != nil {
		response.AbortWithException(c, apiException.UploadFileError, err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			zap.L().Warn("Error closing file", zap.Error(err))
		}
	}(file)

	// 转换到 WebP
	reader, err := objectService.ConvertToWebP(file)
	if errors.Is(err, image.ErrFormat) {
		response.AbortWithException(c, apiException.FileNotImageError, err)
		return
	}
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 上传文件
	objectKey := uuid.NewV1().String() + ".webp"
	err = objectService.SaveObject(reader, objectKey)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	response.JsonSuccessResp(c, gin.H{
		"url": "http://" + c.Request.Host + path.Join("/static", objectKey),
	})
}
