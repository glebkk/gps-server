package model

import "time"

type Movement struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"createdAt"`
}

type MovementCreate struct {
	Latitude  float64 `json:"latitude,string"`
	Longitude float64 `json:"longitude,string"`
}
