package main

import (
	application "Bot/internal"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	app := application.Application{}

	err := godotenv.Load()
	if err != nil {
		return
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Не найдена переменная BOT_TOKEN")
	}

	prefix := os.Getenv("PREFIX")
	if prefix == "" {
		log.Fatal("Не найдена переменная PREFIX")
	}

	application.BotToken = botToken
	application.Prefix = prefix
	app.Start()
}
