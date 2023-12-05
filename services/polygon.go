package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gps_api/db"
	"gps_api/model"
)

type PolygonService struct {
}

func NewPolygonService() *PolygonService {
	return &PolygonService{}
}

func (ps *PolygonService) CreatePolygon(polygon model.PolygonCreate) error {
	_, err := db.Db.Exec(`insert into polygon_areas (geom) values (ST_GeometryFromText($1))`, polygon.Geometry)
	if err != nil {
		fmt.Println(err)
		return errors.New("error create")
	}
	return nil
}

func (ps *PolygonService) GetPolygonByPoint(lat string, long string) (*model.Polygon, error) {
	row := db.Db.QueryRow(`SELECT *
								FROM polygon_areas pa
								where ST_Contains(st_setsrid(pa.geom, 4326), st_setsrid(st_point($1, $2), 4326))
								limit 1`,
		lat, long)

	var pm = &model.Polygon{}
	err := row.Scan(&pm.Id, &pm.Geometry)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error get polygon by movement" + err.Error())
	}

	return pm, nil

}
