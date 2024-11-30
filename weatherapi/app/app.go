package app

import (
	"weatherapi/cache"
	"weatherapi/routes"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	cache.StartRedis()
	defer cache.CloseRedis()

	routes.InitRoutes(r, cache.GetRedis())

	r.Run(":8080")
}
