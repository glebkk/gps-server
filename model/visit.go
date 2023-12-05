package model

import "time"

type Visit struct {
	Id        int        `json:"id"`
	UserId    int        `json:"userId"`
	PolygonId int        `json:"polygonId"`
	TimeEntry time.Time  `json:"timeEntry"`
	TimeExit  *time.Time `json:"timeExit"`
}
