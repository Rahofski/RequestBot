package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
	"fixitpolytech/internal/config"
	"fixitpolytech/internal/handlers"
	"fixitpolytech/internal/services"
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

	err = config.Init()

	if err != nil {
		log.Fatal("Ошибка инициализации конфигурации:", err)
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

	requestService := services.NewRequestService()
    handlers.SetupCommands(bot, requestService)

	bot.Start()

}
