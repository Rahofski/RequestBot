package handlers

import (
	"fixitpolytech/internal/database"
	"fixitpolytech/internal/services"
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func RegisterRequestHandlers(bot *telebot.Bot, requestService *services.RequestService) {
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return DescriptionHandler(c, requestService)
	})

	bot.Handle(&telebot.InlineButton{Unique: "plotnik_job"}, func(c telebot.Context) error {
		return sendJobOptions(c, 0)
	})

	bot.Handle(&telebot.InlineButton{Unique: "electrik_job"}, func(c telebot.Context) error {
		return sendJobOptions(c, 1)
	})

	bot.Handle(&telebot.InlineButton{Unique: "santehnik_job"}, func(c telebot.Context) error {
		return sendJobOptions(c, 2)
	})

	bot.Handle(&telebot.InlineButton{Unique: "skip_photo"}, func(c telebot.Context) error {
		return CompleteRequest(c, requestService)
	})

	bot.Handle(telebot.OnPhoto, func(c telebot.Context) error {
		return HandlePhoto(c, requestService)
	})
}

func DescriptionHandler(c telebot.Context, requestService *services.RequestService) error {
	// Получаем заявку из сессии пользователя
	sessionMutex.Lock()
	req, ok := sessions[c.Sender().ID]
	sessionMutex.Unlock()

	if !ok {
		return c.Send("Ошибка: заявка не найдена")
	}

	req.AdditionalText = c.Text()

	// Обновляем заявку в сессии
	sessionMutex.Lock()
	sessions[c.Sender().ID] = req
	sessionMutex.Unlock()

	keyboard := &telebot.ReplyMarkup{}
	btnPl := keyboard.Data("Плотник", "plotnik_job")
	btnEl := keyboard.Data("Электрик", "electrik_job")
	btnSan := keyboard.Data("Сантехник", "santehnik_job")

	keyboard.Inline(keyboard.Row(btnPl, btnEl, btnSan))

	return c.Send("Выберите, к какому типу относится ваша заявка:", keyboard)
}

func sendJobOptions(c telebot.Context, fieldID int) error {
	sessionMutex.Lock()
	req, ok := sessions[c.Sender().ID]
	sessionMutex.Unlock()

	if !ok {
		return c.Send("Ошибка: заявка не найдена")
	}

	req.FieldID = fieldID
	c.Set("request", req)

	keyboard := &telebot.ReplyMarkup{}
	btnSkip := keyboard.Data("Пропустить", "skip_photo")
	keyboard.Inline(keyboard.Row(btnSkip))

	return c.Send("При желании можете отправить фотографию проблемы (или нажмите 'Пропустить'):", keyboard)
}

func HandlePhoto(c telebot.Context, requestService *services.RequestService) error {
	sessionMutex.Lock()
	req, ok := sessions[c.Sender().ID]
	sessionMutex.Unlock()

	if !ok {
		return c.Send("Ошибка: заявка не найдена")
	}

	photo := c.Message().Photo
	fileURL, err := GetFileURL(c.Bot(), photo.FileID)
	if err != nil {
		log.Printf("Не удалось получить url фото")
	} else {
		req.Photos = append(req.Photos, fileURL)
	}

	return CompleteRequest(c, requestService)
}

func CompleteRequest(c telebot.Context, requestService *services.RequestService) error {
	// Получаем заявку из сессии пользователя
	sessionMutex.Lock()
	req, ok := sessions[c.Sender().ID]
	sessionMutex.Unlock()

	var field string

	if !ok {
		return c.Send("Ошибка: заявка не найдена")
	}

	//Отправка данных на сервер
	if err := database.PostRequest(req); err != nil {
		c.Send("Ошибка: Не удалось отправить заявку на сервер")
	}

	// Создаем заявку в сервисе
	request := requestService.CreateRequest(req.BuildingID, req.FieldID, req.AdditionalText, req.Photos)

	if request.FieldID == 0 {
		field = "Плотник"
	}
	if request.FieldID == 1 {
		field = "Электрик"
	}
	if request.FieldID == 2 {
		field = "Сантехник"
	}

	// Очищаем сессию пользователя
	sessionMutex.Lock()
	delete(sessions, c.Sender().ID)
	sessionMutex.Unlock()

	msg := fmt.Sprintf(
		"Заявка #%d создана!\n"+
			"Здание: %d\n"+
			"URL фото: %s\n"+
			"Сфера: %s\n"+
			"Описание: %s\n"+
			"Статус: %s\n"+
			"Время: %s",
		request.RequestID,
		request.BuildingID,
		request.Photos,
		field,
		request.AdditionalText,
		request.Status,
		request.Time.Format(time.RFC1123),
	)

	return c.Send(msg)
}
