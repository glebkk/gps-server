package model

type Polygon struct {
	Id       int    `json:"id"`
	Geometry string `json:"geometry"`
}

type PolygonCreate struct {
	Geometry string `json:"geometry"`
}
