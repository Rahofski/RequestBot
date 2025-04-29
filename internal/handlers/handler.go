package handlers

import (
	"fixitpolytech/internal/models"
	"fixitpolytech/internal/services"
	_"fixitpolytech/internal/database"

	"log"
	_"strconv"
	"sync"

	"gopkg.in/telebot.v3"
)

var (
    sessions = make(map[int64]*models.Request) // Хранилище сессий (ключ — ID пользователя)
    sessionMutex sync.Mutex                    // Мьютекс для безопасного доступа к хранилищу
    waitingForStatus = make(map[int64]bool)     // Для статуса заявки
)

func SetupCommands(bot *telebot.Bot, requestService *services.RequestService) {
    commands := []telebot.Command{
        {Text: "report", Description: "Оставить заявку"},
        {Text: "help", Description: "Помощь"},
        {Text: "start", Description: "Запуск"},
        {Text: "status", Description: "Получить статус заявки"},
    }

    err := bot.SetCommands(commands)
    if err != nil {
        log.Fatal("Не удалось зарегистрировать команды:", err)
    }

    bot.Handle("/start", StartHandler)
    bot.Handle("/report", ReportHandler)
    bot.Handle("/status", GetStatusHandler)
    bot.Handle("/help", HelpHandler)

    RegisterBuildingHandlers(bot)
    RegisterRequestHandlers(bot, requestService)

    bot.Handle(telebot.OnAddedToGroup, func(c telebot.Context) error {
        c.Send("Бот может быть использован только в личных сообщениях")
        return bot.Leave(c.Chat())
    })
}

func StartHandler(c telebot.Context) error {
    return c.Send("Вас приветствует бот для оставления заявок на проблемы внутри Политехнического университета! Спасибо, что выбрали нас! Для оставления заявки напиши /report")
}

func ReportHandler(c telebot.Context) error {
    keyboard := &telebot.ReplyMarkup{}
    btnStud := keyboard.Data("Учебный корпус", "stud_building")
    btnDorm := keyboard.Data("Общежитие", "dorm_building")

    keyboard.Inline(keyboard.Row(btnStud, btnDorm))

    return c.Send("Выберите тип здания, внутри которого возникла проблема:", keyboard)
}

func GetStatusHandler(c telebot.Context) error {
    sessionMutex.Lock()
    waitingForStatus[c.Sender().ID] = true // Ждём от пользователя номер заявки
    sessionMutex.Unlock()

    return c.Send("Введите номер вашей заявки для получения статуса")
}

func HelpHandler(c telebot.Context) error {
    return c.Send("Сервис по оставлению заявок на проблемы внутри Политеха! Введите команду /report для того, чтобы оставить вашу заявку, /getStatus для получения статуса вашей заявки, /help для получения помощи по работе с ботом")
}