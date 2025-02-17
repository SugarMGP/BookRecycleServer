package main

import (
	"bookrecycle-server/internal/midwares"
	"bookrecycle-server/internal/routes"
	"bookrecycle-server/internal/utils/server"
	"bookrecycle-server/pkg/config"
	"bookrecycle-server/pkg/database"
	"bookrecycle-server/pkg/log"
	"bookrecycle-server/pkg/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	redis.Init()
	routes.Init(r)

	server.Run(r, ":"+config.Config.GetString("server.port"))
}
