package main

import (
	"os"

	"bookrecycle-server/internal/midwares"
	"bookrecycle-server/internal/routes"
	"bookrecycle-server/internal/utils/server"
	"bookrecycle-server/pkg/config"
	"bookrecycle-server/pkg/database"
	"bookrecycle-server/pkg/log"
	"bookrecycle-server/pkg/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 如果配置文件中开启了调试模式
	if !config.Config.GetBool("server.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	log.Init()
	database.Init()
	ws.Init()
	routes.Init(r)

	// 确保 static 目录存在，如果不存在则创建
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		err := os.Mkdir("static", os.ModePerm)
		if err != nil {
			zap.L().Fatal("Failed to create static directory", zap.Error(err))
		}
	}
	r.Static("/static", "./static") // 挂载静态文件目录

	server.Run(r, ":"+config.Config.GetString("server.port"))
}
