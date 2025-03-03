package models

import "time"

type Request struct {
	id          int
	description string
	img         string
	status      string
	building    Building
	time        time.Time
}
