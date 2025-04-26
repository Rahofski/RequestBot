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

	if len(req.Photos) > 0 {
		return c.Send("Вы уже прикрепили фотографию. Нельзя прикрепить больше одной.")
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

	

	// Создаем заявку в сервисе
	request := requestService.CreateRequest(req.BuildingID, req.FieldID, req.AdditionalText, req.Photos)

	build, err := GetBuildingNameByID(request.BuildingID)
	if err != nil {
		return c.Send("Ошибка: Не удалось получить информацию о здании")
	}

	status := request.Status

	if status == "not taken" {
		status = "В обработке"
	}
	if status == "done" {
		status = "Выполнена"
	}
	if status == "in progress" {
		status = "В процессе устранения"
	}

	//Отправка данных на сервер
	requestID, err := database.PostRequest(&request);
	if err != nil {
		c.Send("Ошибка: Не удалось отправить заявку на сервер")
	}

	// Очищаем сессию пользователя
	sessionMutex.Lock()
	delete(sessions, c.Sender().ID)
	sessionMutex.Unlock()



	msg := fmt.Sprintf(
		"Заявка создана!\n"+
			"ID заявки: %d\n"+
			"Здание: %s\n"+
			"Описание: %s\n"+
			"Статус: %s\n"+
			"Время: %s",
		requestID,
		build,
		request.AdditionalText,
		status,
		request.Time.Format(time.RFC1123),
	)

	return c.Send(msg)
}
