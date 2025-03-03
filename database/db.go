package database

import (
    models "fixitpolytech/internal/models"
)
//метод для получения всех корпусовэ(нужно реализовать через подключение к бд)

func GetAllBuildings() []models.Building{
    return []models.Building{
        {Name: "11 корпус", Address: "улица Обручевых, 1", BldType: "stud"},
        {Name: "9 корпус", Address: "Политехническая улица, 21", BldType: "stud"},
        {Name: "6 общага", Address: "Улица Харченко, 16", BldType: "dorm"},
    }
}


