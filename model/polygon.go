package model

import (
	"github.com/paulmach/orb/geojson"
)

type Polygon struct {
	Id       int              `json:"id"`
	Geometry geojson.Geometry `json:"geometry"`
}

type PolygonCreate struct {
	Geometry geojson.Geometry `json:"geometry"`
}
