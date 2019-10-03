package main

import (
	"fmt"
	"log"
	"os"

	owm "github.com/briandowns/openweathermap"
	"github.com/yanzay/tbot/v2"
)

type weatherBot struct {
	client       *tbot.Client
	weatherAPI   *owm.CurrentWeatherData
	weatherIcons map[string]string
}

func getWeatherIcons() map[string]string {
	return map[string]string{
		"Thunderstorm": "☁☔⚡",
		"Drizzle":      "☔",
		"Rain":         "☁☔",
		"Snow":         "❄",
		"Clear":        "☀",
		"Clouds":       "⛅",
	}
}

func (wb *weatherBot) StartHandler(msg *tbot.Message) {
	wb.client.SendMessage(msg.Chat.ID, "Напишите мне название вашего города и я скажу вам температуру")
}

func (wb *weatherBot) getWeatherIcon() string {
	if len(wb.weatherAPI.Weather) != 0 {
		if icon, exist := wb.weatherIcons[wb.weatherAPI.Weather[0].Main]; exist {
			return icon
		}
	}
	return ""
}

func (wb *weatherBot) WeatherHandler(msg *tbot.Message) {
	err := wb.weatherAPI.CurrentByName(msg.Text)
	if err != nil {
		wb.client.SendMessage(msg.Chat.ID, "Не могу найти ваш город")
		return
	}
	weather := wb.getWeatherIcon()
	wb.client.SendMessage(msg.Chat.ID, fmt.Sprintf("В городе %s %.2f С° %s", msg.Text, wb.weatherAPI.Main.Temp, weather))
}

func Client(bot *tbot.Server) (*weatherBot, error) {
	wAPI, err := owm.NewCurrent("C", "en", os.Getenv("OWM_KEY"))
	if err != nil {
		return nil, err
	}
	return &weatherBot{
		client:       bot.Client(),
		weatherAPI:   wAPI,
		weatherIcons: getWeatherIcons(),
	}, nil
}

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_KEY"))
	client, err := Client(bot)
	if err != nil {
		log.Fatal(err)
	}
	bot.HandleMessage(`/start`, client.StartHandler)
	bot.HandleMessage(`.*`, client.WeatherHandler)
	err = bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
