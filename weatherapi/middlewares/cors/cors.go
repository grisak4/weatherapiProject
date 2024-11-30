package cors

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetCors(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
