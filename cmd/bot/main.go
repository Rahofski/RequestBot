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
	// Try to load .env file from multiple locations
	envLocations := []string{
		"/app/.env",    // Docker container location
		"./.env",       // Local development
		"../../.env",   // Alternative local path
	}
	
	var envLoaded bool
	for _, loc := range envLocations {
		if err := godotenv.Load(loc); err == nil {
			envLoaded = true
			break
		}
	}
	
	if !envLoaded {
		log.Println("Warning: No .env file found, relying on environment variables")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN not found in environment variables!")
	}

	err := config.Init()
	if err != nil {
		log.Fatal("Configuration initialization error:", err)
	}

	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	// Create bot
	bot, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	requestService := services.NewRequestService()
	handlers.SetupCommands(bot, requestService)

	bot.Start()
}