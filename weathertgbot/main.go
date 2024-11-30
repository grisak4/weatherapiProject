package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	botToken = "8102204031:AAFXpIZurtopq3e77l-ZQJKAmoVIwwIaIIU"
	api_url  = "http://localhost:8080"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	log.Println("[BOT] Bot has started!")
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Printf("[HANDLER] Got a message. || User: %s ;", update.Message.From.Username)

	if update.Message.Location != nil {
		handleLocation(ctx, b, update)
		return
	}

	switch update.Message.Text {
	case "/start":
		cmdStart(ctx, b, update)
	case "/weather":
		cmdWeatherInfo(ctx, b, update)
	}
}

func cmdStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет! Это погодный бот. Отправь мне свою геолокацию, чтобы получить прогноз погоды.",
	})
}

func cmdWeatherInfo(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отправь мне свою геолокацию, чтобы узнать прогноз погоды!",
	})
}

func handleLocation(ctx context.Context, b *bot.Bot, update *models.Update) {
	latitude := update.Message.Location.Latitude
	longitude := update.Message.Location.Longitude

	url := fmt.Sprintf("%s/%f/%f", api_url, latitude, longitude)

	resp, err := GetBody(url)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось получить прогноз погоды. Попробуйте позже.",
		})
		return
	}

	res := editJson(resp)
	message := editMap(res)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: "Markdown",
	})
}

func GetBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func editJson(jsonData []byte) map[string]float64 {
	var outer struct {
		Response string `json:"response"`
	}

	err := json.Unmarshal(jsonData, &outer)
	if err != nil {
		fmt.Println("Ошибка декодирования внешнего JSON:", err)
		return nil
	}

	var inner struct {
		Hourly struct {
			Time          []string  `json:"time"`
			Temperature2m []float64 `json:"temperature_2m"`
		} `json:"hourly"`
	}

	err = json.Unmarshal([]byte(outer.Response), &inner)
	if err != nil {
		fmt.Println("Ошибка декодирования внутреннего JSON:", err)
		return nil
	}

	res := make(map[string]float64)
	for i := 0; i < len(inner.Hourly.Time); i++ {
		res[inner.Hourly.Time[i]] = inner.Hourly.Temperature2m[i]
	}

	return res
}

func editMap(hashTab map[string]float64) string {
	type TimeTemp struct {
		Time string
		Temp float64
	}
	var data []TimeTemp

	for timeStr, temp := range hashTab {
		parsedTime, err := time.Parse("2006-01-02T15:04", timeStr)
		if err != nil {
			continue
		}

		hour := parsedTime.Hour()
		if hour >= 0 && hour <= 23 {
			data = append(data, TimeTemp{Time: timeStr, Temp: temp})
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Time < data[j].Time
	})

	result := "🌤️ *Прогноз погоды на день:*\n\n"
	for _, entry := range data {
		parsedTime, _ := time.Parse("2006-01-02T15:04", entry.Time)
		timeStr := parsedTime.Format("15:04")

		var emoji string
		switch {
		case entry.Temp <= -10:
			emoji = "❄️"
		case entry.Temp <= 0:
			emoji = "🧊"
		case entry.Temp <= 10:
			emoji = "🌬️"
		case entry.Temp <= 20:
			emoji = "☁️"
		default:
			emoji = "☀️"
		}

		result += fmt.Sprintf("%s %s: %.1f°C\n", emoji, timeStr, entry.Temp)
	}

	return result
}
