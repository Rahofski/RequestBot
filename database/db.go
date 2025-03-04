package database

import (
    models "fixitpolytech/internal/models"
)
//метод для получения всех корпусовэ(нужно реализовать через подключение к бд)

func GetAllBuildings() []models.Building{
    return []models.Building{
        {Name: "11 корпус", EngName: "11corp", Address: "улица Обручевых, 1", BldType: "stud"},
        {Name: "9 корпус", EngName: "9corp", Address: "Политехническая улица, 21", BldType: "stud"},
        {Name: "6 общага", EngName: "6dorm", Address: "Улица Харченко, 16", BldType: "dorm"},
        {Name: "3 общага", EngName: "3dorm", Address: "Лесной проспект, 65к3", BldType: "dorm"},
    }
}


