package database

import (
	models "fixitpolytech/internal/models"
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
func GetAllBuildings() []models.Building {
	return []models.Building{
		{BuildingID:0, Name: "11 корпус", Address: "улица Обручевых, 1", BldType: "stud"},
		{BuildingID:1, Name: "9 корпус",  Address: "Политехническая улица, 21", BldType: "stud"},
		{BuildingID:2, Name: "6 общага",  Address: "Улица Харченко, 16", BldType: "dorm"},
		{BuildingID:3, Name: "3 общага",  Address: "Лесной проспект, 65к3", BldType: "dorm"},
	}
}
