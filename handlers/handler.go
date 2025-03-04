package handlers

import (
	"fixitpolytech/internal/database"
	"fixitpolytech/internal/models"
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

var RequestList []models.Request

func SetupCommands(bot *telebot.Bot) {
	commands := []telebot.Command{
		{Text: "report", Description: "Оставить заявку"},
		{Text: "help", Description: "Помощь"},
		{Text: "start", Description: "Запуск"},
	}

	err := bot.SetCommands(commands)
	if err != nil {
		log.Fatal("Не удалось зарегистрировать команды:", err)
	}
	bot.Handle("/start", StartHandler)
	bot.Handle("/report", ReportHandler)
	bot.Handle("/help", HelpHandler)

	bot.Handle(&telebot.InlineButton{Unique: "stud_building"}, func(c telebot.Context) error {
		return sendBuildingOptions(c, "stud")
	})
	
	bot.Handle(&telebot.InlineButton{Unique: "dorm_building"}, func(c telebot.Context) error {
		return sendBuildingOptions(c, "dorm")
	})

	registerBuildingHandlers(bot)

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

	return c.Send("Выберите тип здания, внутри которого возникла проблема:", keyboard)
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
				uniqueID := "bld_" + building.EngName
				btn := keyboard.Data(building.Name, uniqueID)
				log.Printf("Создаём кнопку с Unique: %s", uniqueID)
				rows = append(rows, keyboard.Row(btn))
		}
	}

	keyboard.Inline(rows...)

	_ = c.Respond()
	return c.Send("Выберите конкретное здание:", keyboard)
}

func registerBuildingHandlers(bot *telebot.Bot) {
	buildings := database.GetAllBuildings()

	for _, building := range buildings {
		b:= building

		uniqueID := "bld_" + building.EngName
		log.Printf("Регистрируем обработчик с Unique: %s", uniqueID)
		bot.Handle(&telebot.InlineButton{Unique: uniqueID}, func(c telebot.Context) error {
			return ParseBldButton(c, b)
		})
	}
}

func ParseBldButton(c telebot.Context, bld models.Building) error {
	
	req := models.Request{
		Building: bld,
	}
	RequestList = append(RequestList, req)
	
	err := c.Send("Опишите проблему (уточните местоположение внутри здания, дайте дополнительную информацию и т.д.)")
	if err != nil {
		return err
	}
	
	// Регистрируем временный обработчик для следующего сообщения
	c.Bot().Handle(telebot.OnText, func(ctx telebot.Context) error {
		return DescriptionHandler(ctx)
	})
	
	_ = c.Respond()
	return nil
}

func HelpHandler(c telebot.Context) error {
	return c.Send("Сервис по оставлению заявок на проблемы внутри Политеха! Введите команду /report для того, чтобы оставить вашу заявку")
}

func StartHandler(c telebot.Context) error {
	return c.Send("Вас приветствует бот для оставления заявок на проблемы внутри Политехнического университета! Спасибо, что выбрали нас! Для оставления заявки напиши /report")
}

func DescriptionHandler(c telebot.Context) error {
	if len(RequestList) == 0 {
		return c.Send("Ошибка: заявка не найдена")
	}
	
	lastIndex := len(RequestList) - 1
	RequestList[lastIndex].Description = c.Text()
	
	keyboard := &telebot.ReplyMarkup{}
	btnSkip := keyboard.Data("Пропустить", "skip_photo")
	keyboard.Inline(keyboard.Row(btnSkip))
	
	err := c.Send("При желании можете отправить фотографию проблемы (или нажмите 'Пропустить'):", keyboard)
	if err != nil {
		return err
	}
	
	c.Bot().Handle(&telebot.InlineButton{Unique: "skip_photo"}, func(ctx telebot.Context) error {
		return CompleteRequest(ctx, lastIndex)
	})
	
	c.Bot().Handle(telebot.OnPhoto, func(ctx telebot.Context) error {
		photo := ctx.Message().Photo
		RequestList[lastIndex].Img, err = GetFileURL(c.Bot(), photo.FileID)
		if err != nil{
			log.Printf("Не удалось получить url фото")
		}
		return CompleteRequest(ctx, lastIndex)
	})
	
	return nil
}

func CompleteRequest(c telebot.Context, index int) error {
	RequestList[index].Status = "in process"
	RequestList[index].Time = time.Now()
	RequestList[index].Id = index + 1 // Простая генерация ID
	
	msg := fmt.Sprintf(
		"Заявка #%d создана!\n"+
			"Здание: %s\n"+
			"URL фото: %s\n"+
			"Описание: %s\n"+
			"Статус: %s\n"+
			"Время: %s",
		RequestList[index].Id,
		RequestList[index].Building.Name,
		RequestList[index].Img,
		RequestList[index].Description,
		RequestList[index].Status,
		RequestList[index].Time.Format(time.RFC1123),
	)
	
	c.Bot().Handle(telebot.OnText, func(ctx telebot.Context) error {
        return nil
    })
    c.Bot().Handle(telebot.OnPhoto, func(ctx telebot.Context) error {
        return nil
    })

	//здесь должна быть отправка заявка в БД
	
	_ = c.Respond()
	return c.Send(msg)
}

