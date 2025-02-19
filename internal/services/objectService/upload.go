package objectService

import (
	"bytes"
	"image"
	_ "image/gif" // 注册解码器
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
	_ "golang.org/x/image/bmp" // 注册解码器
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// SizeLimit 上传大小限制
const SizeLimit = humanize.MByte * 10

// SaveObject 根据 ObjectKey 保存文件
func SaveObject(reader io.Reader, objectKey string) error {
	// 根据 objectKey 解析出文件的路径
	relativePath := filepath.Join("static", objectKey)

	// 创建文件
	outFile, err := os.Create(relativePath)
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			zap.L().Warn("Error closing file", zap.Error(err))
		}
	}(outFile)

	// 写入文件
	_, err = io.Copy(outFile, reader)
	return err
}

// ConvertToWebP 将图片转换为 WebP 格式
func ConvertToWebP(reader io.Reader) (io.Reader, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = webp.Encode(&buf, img, &webp.Options{Quality: 100})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}
