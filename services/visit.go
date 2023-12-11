package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gps_api/db"
	"gps_api/model"
	"time"
)

type VisitService struct {
}

func NewVisitService() *VisitService {
	return &VisitService{}
}

func (vs *VisitService) CheckVisitOpen(userId int) (bool, error) {
	row := db.Db.QueryRow(`SELECT time_exit is null
							FROM visits vs
							where vs.user_id = $1
							order by time_entry desc
							limit 1`,
		userId)

	var visitIsOpen bool
	err := row.Scan(&visitIsOpen)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		fmt.Println("err in check is visit open", err)
		return false, errors.New("err in check is visit open" + err.Error())
	}

	return visitIsOpen, nil
}

func (vs *VisitService) CloseVisit(userId int) error {
	_, err := db.Db.Exec(`update visits
									set time_exit = now()
									where id = (
										SELECT id
										FROM visits v 
										where user_id = $1 and time_exit is null
										limit 1
									)`,
		userId)
	if err != nil {
		fmt.Println("err in close visit", err)
		return errors.New("err in close visit" + err.Error())
	}
	fmt.Println("close visit", time.Now().Format(time.TimeOnly))

	return nil
}

func (vs *VisitService) OpenVisit(userId int, polygonId int) error {
	_, err := db.Db.Exec(`insert into visits (user_id, polygon_id) values ($1, $2)`,
		userId, polygonId)
	if err != nil {
		fmt.Println("err in open visit", err)
		return errors.New("err in open visit" + err.Error())
	}
	fmt.Println("open visit", time.Now().Format(time.TimeOnly))

	return nil
}

func (vs *VisitService) GetLastUserVisits(userId int) ([]model.Visit, error) {
	rows, err := db.Db.Query(`SELECT t1.*
									FROM visits t1
									INNER JOIN (
										SELECT polygon_id, MAX(time_entry) AS max_time
										FROM visits
										WHERE user_id = $1
										GROUP BY polygon_id
									) t2 ON t1.polygon_id = t2.polygon_id AND t1.time_entry = t2.max_time
									WHERE t1.user_id = $1;`, userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	var visits []model.Visit
	var visit model.Visit

	for rows.Next() {
		err := rows.Scan(&visit.Id, &visit.UserId, &visit.PolygonId, &visit.TimeEntry, &visit.TimeExit)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		visits = append(visits, visit)
	}

	return visits, nil
}
