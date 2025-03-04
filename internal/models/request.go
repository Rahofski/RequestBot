package models

import "time"

type Request struct {
	Id          int
	Description string
	Img         string
	Status      string
	Building    Building
	Time        time.Time
}
