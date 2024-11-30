package routes

import (
	"weatherapi/services/weather"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func InitRoutes(router *gin.Engine, rdb *redis.Client) {
	router.GET("/:lat/:lon", func(ctx *gin.Context) {
		weather.GetData(ctx, rdb)
	})
}
