package model

type User struct {
	Id   int    `json:"id"`
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}
