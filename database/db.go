package database

//метод для получения всех корпусовэ(нужно реализовать через подключение к бд)


type Building struct{
    name string
    adress string
}

func GetAllBuildings() []Building{
    return []Building{{"11 корпус", "улица Обручевых, 1"}, {"9 корпус", "Политехническая улица, 21"}}
}

//метод для отправки заявки

