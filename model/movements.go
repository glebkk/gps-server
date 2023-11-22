package model

import "time"

type Movement struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	CreatedAt time.Time `json:"createdAt"`
}

type MovementCreate struct {
	Latitude  float32 `json:"latitude,string"`
	Longitude float32 `json:"longitude,string"`
}
