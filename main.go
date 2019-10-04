package main

import (
	"log"
	"os"

	"github.com/Millfort/weatherbot/weatherbot"
)

func main() {
	bot, err := weatherbot.New(os.Getenv("TELEGRAM_KEY"), os.Getenv("OWM_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	bot.HandleMessage(`/start`, bot.StartHandler)
	bot.HandleMessage(`.*`, bot.WeatherHandler)
	err = bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
