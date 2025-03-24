package handlers

import (
	"fixitpolytech/internal/database"
	"fixitpolytech/internal/services"
	"fixitpolytech/internal/services/gigachat"
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func RegisterRequestHandlers(bot *telebot.Bot, requestService *services.RequestService) {
    bot.Handle(telebot.OnText, func(c telebot.Context) error {
        sessionMutex.Lock()
        _, ok := sessions[c.Sender().ID]
        sessionMutex.Unlock()

        if !ok {
            // Если активной заявки нет, игнорируем сообщение
            return nil
        }

        return DescriptionHandler(c, requestService)
    })

    bot.Handle(&telebot.InlineButton{Unique: "skip_photo"}, func(c telebot.Context) error {
        return CompleteRequest(c, requestService)
    })

    bot.Handle(telebot.OnPhoto, func(c telebot.Context) error {
		sessionMutex.Lock()
        _, ok := sessions[c.Sender().ID]
        sessionMutex.Unlock()

        if !ok {
            // Если активной заявки нет, игнорируем сообщение
            return nil
        }
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

	accessToken := gigachat.GetAccessToken()

	answer, err := gigachat.CheckIfValid(req.AdditionalText, accessToken)
	if err != nil {
		return c.Send("Ошибка: Не удалось проверить валидность заявки")
	}
	// Проверяем валидность заявки
	if !answer {
		// Если заявка невалидна — обнуляем её
		sessionMutex.Lock()
		delete(sessions, c.Sender().ID)
		sessionMutex.Unlock()
		return c.Send("Ошибка: заявка невалидна, попробуйте снова.")
	}

	field, err := gigachat.DefineField(req.AdditionalText, accessToken)

	if err != nil {
		return c.Send("Ошибка: Не удалось определить сферу заявки")
	}

	if field == -1 {
		return c.Send("Ошибка: Не удалось определить сферу заявки")
	}

	req.FieldID = field

	// Обновляем заявку в сессии
	sessionMutex.Lock()
	sessions[c.Sender().ID] = req
	sessionMutex.Unlock()

	// Кнопка для пропуска фото
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

	if !ok {
		return c.Send("Ошибка: заявка не найдена")
	}

	//Отправка данных на сервер
	if err := database.PostRequest(req); err != nil {
		c.Send("Ошибка: Не удалось отправить заявку на сервер")
	}

	// Создаем заявку в сервисе
	request := requestService.CreateRequest(req.BuildingID, req.FieldID, req.AdditionalText, req.Photos)

	// Очищаем сессию пользователя
	sessionMutex.Lock()
	delete(sessions, c.Sender().ID)
	sessionMutex.Unlock()

	msg := fmt.Sprintf(
		"Заявка #%d создана!\n"+
			"Здание: %d\n"+
			"URL фото: %s\n"+
			"Описание: %s\n"+
			"Статус: %s\n"+
			"Время: %s",
		request.RequestID,
		request.BuildingID,
		request.Photos,
		request.AdditionalText,
		request.Status,
		request.Time.Format(time.RFC1123),
	)

	return c.Send(msg)
}
