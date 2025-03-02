package handlers

import (
	"log"

	"gopkg.in/telebot.v3"
)

func SetupCommands(bot *telebot.Bot) {
	commands := []telebot.Command{
		{Text: "report", Description: "Оставить заявку"},
		{Text: "help", Description: "Помощь"},
	}

    err := bot.SetCommands(commands)
	if err != nil {
		log.Fatal("Не удалось зарегистрировать команды:", err)
	}

	bot.Handle("/report", ReportHandler)

	bot.Handle("/help", HelpHandler)
}

func ReportHandler(c telebot.Context) error {
    c.Send("С")
    
    
    return nil
}

func HelpHandler(c telebot.Context) error {
    return c.Send("Помогу с проблемами Политеха! Введите команду/report для того, чтобы оставить вашу заявку")
}
