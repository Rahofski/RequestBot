package database

import (
	"bytes"
	"encoding/json"
	models "fixitpolytech/internal/models"
	"fmt"
	"log"
	"net/http"
)

//метод для получения всех корпусов(нужно реализовать через подключение к бд)
/*func GetAllBuildings() []models.Building {
	response, err := http.Get("http://localhost:3000/api/buildings")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var buildings []Building
	if err := json.Unmarshal(body, &buildings); err != nil {
		log.Fatal(err)
	}
	return buildings
}*/

func PostRequest(req *models.Request) error {
	// Преобразуем запрос в JSON
	payload, err := json.Marshal(*req)
	if err != nil {
		log.Println("Ошибка при преобразовании данных в JSON:", err)
		return err
	}
	//Создаем POST-запрос
	response, err := http.Post("http://localhost:3000/api/requests", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return err
	}
	defer response.Body.Close()
	//Проверяем ответ сервера
	if response.StatusCode != http.StatusOK {
		log.Printf("Сервер вернул статус: %s\n", response.Status)
		return fmt.Errorf("Ошибка сервера: %s", response.Status)
	}

	return nil
}
func GetAllBuildings() []models.Building {
	return []models.Building{
		{BuildingID: 0, Name: "11 корпус", Address: "улица Обручевых, 1", BldType: "stud"},
		{BuildingID: 1, Name: "9 корпус", Address: "Политехническая улица, 21", BldType: "stud"},
		{BuildingID: 2, Name: "6 общага", Address: "Улица Харченко, 16", BldType: "dorm"},
		{BuildingID: 3, Name: "3 общага", Address: "Лесной проспект, 65к3", BldType: "dorm"},
	}
}
