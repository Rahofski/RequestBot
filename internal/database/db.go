package database

import (
	"bytes"
	"encoding/json"
	models "fixitpolytech/internal/models"
	"fmt"
	 "io"
	"log"
	"net/http"
)

//метод для получения всех корпусов(нужно реализовать через подключение к бд)
func GetAllBuildings() []models.Building {
	response, err := http.Get("http://localhost:8080/api/buildings")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var buildings []models.Building
	if err := json.Unmarshal(body, &buildings); err != nil {
		log.Fatal("error while parsing buildings: ", err)
	}
	return buildings
}

func PostRequest(req *models.Request) (int, error) {
	// Преобразуем запрос в JSON
	payload, err := json.Marshal(*req)
	if err != nil {
		log.Println("Ошибка при преобразовании данных в JSON:", err)
		return 0, err
	}
	//Создаем POST-запрос
	response, err := http.Post("http://localhost:8080/api/request/add", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return 0, err
	}
	defer response.Body.Close()

	//Проверяем ответ сервера
	if response.StatusCode != http.StatusOK {
		log.Printf("Сервер вернул статус: %s\n", response.Status)
		return 0, fmt.Errorf("Ошибка сервера: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseData models.RequestResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Fatal("error while parsing response: ", err)
	}

	return responseData.RequestID, nil
}

func GetRequestStatus(requestID int) (string, error) {
	response, err := http.Get("http://localhost:8080/api/status/" + fmt.Sprint(requestID))
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Сервер вернул статус: %s\n", response.Status)
		return "", fmt.Errorf("Ошибка сервера: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseData models.StatusResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Fatal("error while parsing response: ", err)
	} 

	return responseData.Status, nil

}

// func GetAllBuildings() []models.Building {
// 	return []models.Building{
// 		{BuildingID: 0, Name: "11 корпус", Address: "улица Обручевых, 1", BldType: "корпус"},
// 		{BuildingID: 1, Name: "9 корпус", Address: "Политехническая улица, 21", BldType: "корпус"},
// 		{BuildingID: 2, Name: "6 общага", Address: "Улица Харченко, 16", BldType: "общежитие"},
// 		{BuildingID: 3, Name: "3 общага", Address: "Лесной проспект, 65к3", BldType: "общежитие"},
// 	}
// }
