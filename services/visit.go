package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gps_api/db"
)

type VisitService struct {
}

func NewVisitService() *VisitService {
	return &VisitService{}
}

func (vs *VisitService) CheckVisitOpen(userId int) (bool, error) {
	row := db.Db.QueryRow(`SELECT time_exit is not null
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
		fmt.Println(err)
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
		fmt.Println(err)
		return errors.New("err in close visit" + err.Error())
	}

	return nil
}

func (vs *VisitService) OpenVisit(userId int, polygonId int) error {
	_, err := db.Db.Exec(`insert into visits (user_id, polygon_id) values ($1, $2)`,
		userId, polygonId)
	if err != nil {
		fmt.Println(err)
		return errors.New("err in open visit" + err.Error())
	}

	return nil
}
