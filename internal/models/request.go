package models

import "time"

// type Request struct {
// 	Id          int
// 	Description string
// 	Img         string
// 	Status      string
// 	Building    Building
// 	Time        time.Time
// }

type Request struct {
	RequestID      int      `json:"request_id"`
	BuildingID     int      `json:"building_id"`
	FieldID        int      `json:"field_id"`
	AdditionalText string   `json:"additional_text"`
	Status         string   `json:"status"`
	Photos         string `json:"photos"`
	Time		   time.Time `json:"time"`
}

type RequestResponse struct {
	RequestID      int      `json:"request_id"`
}


//ID по хорошему надо убрать и здесь и в handler, потому что бд сама задает ID
//Нужно посмотреть как он будет форматы time совмещать (формат бд и формат go)
