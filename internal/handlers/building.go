package handlers

import (
    "fixitpolytech/internal/database"
    "fixitpolytech/internal/models"
    "fmt"
    "log"

    "gopkg.in/telebot.v3"
)

func RegisterBuildingHandlers(bot *telebot.Bot) {
    // Обработчики для статических кнопок (учебный корпус и общежитие)
    bot.Handle(&telebot.InlineButton{Unique: "stud_building"}, func(c telebot.Context) error {
        return sendBuildingOptions(c, "корпус")
    })

    bot.Handle(&telebot.InlineButton{Unique: "dorm_building"}, func(c telebot.Context) error {
        return sendBuildingOptions(c, "общежитие")
    })

    // Обработчики для динамических кнопок зданий
    buildings := database.GetAllBuildings()
    for _, building := range buildings {
        b := building
        uniqueID := "bld_" + fmt.Sprint(b.BuildingID)
        log.Printf("Регистрируем обработчик с Unique: %s", uniqueID)
        bot.Handle(&telebot.InlineButton{Unique: uniqueID}, func(c telebot.Context) error {
            return ParseBldButton(c, b)
        })
    }
}

func sendBuildingOptions(c telebot.Context, bldType string) error {
    buildings := database.GetAllBuildings()

    if len(buildings) == 0 {
        return c.Send("Зданий не найдено.")
    }

    keyboard := &telebot.ReplyMarkup{}
    var rows []telebot.Row

    for _, building := range buildings {
        if building.BldType == bldType {
            uniqueID := "bld_" + fmt.Sprint(building.BuildingID)
            btn := keyboard.Data(building.Name, uniqueID)
            log.Printf("Создаём кнопку с Unique: %s", uniqueID)
            rows = append(rows, keyboard.Row(btn))
        }
    }

    keyboard.Inline(rows...)

    _ = c.Respond()
    return c.Send("Выберите конкретное здание:", keyboard)
}

func GetBuildingNameByID(buildingID int) (string, error) {
    buildings := database.GetAllBuildings() // Получаем список всех зданий

    for _, building := range buildings {
        if building.BuildingID == buildingID {
            return building.Name, nil // Нашли здание, возвращаем его имя
        }
    }

    return "", fmt.Errorf("здание с ID %d не найдено", buildingID) // Если здание не найдено
}

func ParseBldButton(c telebot.Context, bld models.Building) error {
    req := models.Request{
        BuildingID: bld.BuildingID,
    }

    // Сохраняем заявку в сессию пользователя
    sessionMutex.Lock()
    sessions[c.Sender().ID] = &req
    sessionMutex.Unlock()

    return c.Send("Опишите проблему (уточните местоположение внутри здания, предоставьте дополнительную информацию и т.д.). Пожалуйста, не пишите ничего лишнего, иначе бот не примет вашу заявку.")
}