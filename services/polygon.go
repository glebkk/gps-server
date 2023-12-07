package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"gps_api/db"
	"gps_api/model"
)

type PolygonService struct {
}

func NewPolygonService() *PolygonService {
	return &PolygonService{}
}

func (ps *PolygonService) CreatePolygon(polygon model.PolygonCreate) error {
	geojsonRaw, err := json.Marshal(polygon.Geometry)
	_, err = db.Db.Exec(`insert into polygon_areas (geom) values (ST_FlipCoordinates(st_geomfromgeojson($1)))`, geojsonRaw)
	if err != nil {
		fmt.Println(err)
		return errors.New("error create")
	}
	return nil
}

func (ps *PolygonService) GetAll() ([]model.Polygon, error) {

	var polygons []model.Polygon
	rows, err := db.Db.Query("select id, st_asgeojson(ST_FlipCoordinates( pa.geom )) from polygon_areas pa ")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var geometryStr string

		err := rows.Scan(&id, &geometryStr)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		var geometry geojson.Geometry
		err = json.Unmarshal([]byte(geometryStr), &geometry)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		polygon := model.Polygon{
			Id:       id,
			Geometry: geometry,
		}

		polygons = append(polygons, polygon)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return polygons, nil
}

func (ps *PolygonService) GetPolygonByPoint(long string, lat string) (*model.Polygon, error) {
	row := db.Db.QueryRow(`SELECT id, st_asgeojson(ST_FlipCoordinates( pa.geom ))
								FROM polygon_areas pa
								where ST_Contains(st_setsrid(pa.geom, 4326), st_setsrid(st_point($1, $2), 4326))
								limit 1`,
		long, lat)

	var id int
	var geometryStr string

	err := row.Scan(&id, &geometryStr)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		fmt.Println("get polygon by point", err)
		return nil, err
	}

	fmt.Println(geometryStr)
	var geometry geojson.Geometry
	err = json.Unmarshal([]byte(geometryStr), &geometry)
	if err != nil {
		fmt.Println("unmarshal err", err)
		return nil, err
	}

	polygon := &model.Polygon{
		Id:       id,
		Geometry: geometry,
	}

	if err != nil {
		fmt.Println("error get polygon by movement", err)
		return nil, errors.New("error get polygon by movement" + err.Error())
	}

	return polygon, nil
}
