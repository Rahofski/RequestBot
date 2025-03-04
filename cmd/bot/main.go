package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
	"github.com/joho/godotenv"

	"fixitpolytech/internal/handlers"
)

func main() {
	// Загружаем переменные окружения из .env
	err := godotenv.Load("../../.env") // Указываем путь к файлу .env из cmd/bot
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	token := os.Getenv("TOKEN")

	if token == "" {
		log.Fatal("TOKEN не найден!")
	}

	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	// Создаём бота
	bot, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	handlers.SetupCommands(bot)

	bot.Start()
}
