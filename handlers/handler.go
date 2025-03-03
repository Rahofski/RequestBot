package handlers

import (
	"fixitpolytech/internal/database"
	"fixitpolytech/internal/models"
	"log"
	"gopkg.in/telebot.v3"
)

var RequestList []models.Request

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

	bot.Handle(&telebot.InlineButton{Unique: "stud_building"}, func(c telebot.Context) error {
		return sendBuildingOptions(c, "stud")
	})
	bot.Handle(&telebot.InlineButton{Unique: "dorm_building"}, func(c telebot.Context) error {
		return sendBuildingOptions(c, "dorm")
	})

	bot.Handle(telebot.OnAddedToGroup, func(c telebot.Context) error {
		c.Send("Бот может быть использован только в личных сообщениях")
		return bot.Leave(c.Chat()) 
	})
	
}

// Выбор типа здания
func ReportHandler(c telebot.Context) error {
	keyboard := &telebot.ReplyMarkup{}
	btnStud := keyboard.Data("Учебный корпус", "stud_building")
	btnDorm := keyboard.Data("Общежитие", "dorm_building")

	keyboard.Inline(
		keyboard.Row(btnStud, btnDorm),
	)

	c.Send("Выберите тип здания, внутри которого возникла проблема:", keyboard)


	return c.Send("Ваша заявка принята!")
}

// Отправка списка зданий в зависимости от типа
func sendBuildingOptions(c telebot.Context, bldType string) error {
	buildings := database.GetAllBuildings()

	if len(buildings) == 0 {
		return c.Send("Зданий не найдено.")
	}

	keyboard := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for _, building := range buildings {
		if building.BldType == bldType {
			btn := keyboard.Data(building.Name, "bld_"+building.Name)
			rows = append(rows, keyboard.Row(btn)) // Добавляем в строки клавиатуры
		}
	}

	keyboard.Inline(rows...)

	_ = c.Respond()
	return c.Send("Выберите конкретное здание:", keyboard)
}

func HelpHandler(c telebot.Context) error {
	return c.Send("Сервис по оставлению заявок на проблемы внутри Политеха! Введите команду /report для того, чтобы оставить вашу заявку")
}
