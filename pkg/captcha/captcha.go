package captcha

import (
	"time"

	"github.com/patrickmn/go-cache"
	images "github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
	"go.uber.org/zap"
)

// SlideCaptcha 滑块验证码生成器
var SlideCaptcha slide.Captcha

// CaptchaStore 验证码缓存存储
var CaptchaStore *cache.Cache

// Init 初始化
func Init() {
	builder := slide.NewBuilder()

	imgs, err := images.GetImages()
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	graphs, err := getTileGraph()
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	builder.SetResources(
		slide.WithGraphImages(graphs),
		slide.WithBackgrounds(imgs),
	)

	SlideCaptcha = builder.Make()
	CaptchaStore = cache.New(5*time.Minute, 10*time.Minute)
}

func getTileGraph() ([]*slide.GraphImage, error) {
	graphs, err := tiles.GetTiles()
	if err != nil {
		return nil, err
	}

	newGraphs := make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}
	return newGraphs, nil
}
