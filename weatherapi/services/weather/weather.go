package weather

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"weatherapi/util/requests"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// 53.9059, 86.719
func GetData(c *gin.Context, rdb *redis.Client) {
	lat := c.Param("lat")
	lon := c.Param("lon")

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=temperature_2m&timezone=Europe/Moscow&forecast_days=1",
		lat, lon)

	key := fmt.Sprintf("forecast %s:%s", lat, lon)

	if data, err := rdb.Get(c, key).Result(); err != nil {
		log.Println("[REDIS] Empty key")
		log.Println("[REDIS.INFO] ", data)

		body, err := requests.GetBody(url)
		if err != nil {
			log.Printf("Ошибка при запросе: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при запросе к внешнему API"})
			return
		}

		rdb.Set(c, key, string(body), time.Minute*1)

		c.JSON(200, gin.H{
			"response": string(body),
		})
	} else {
		c.JSON(200, gin.H{
			"response": string(data),
		})
	}
}
