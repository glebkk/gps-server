package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gps_api/db"
	"time"
)

type VisitService struct {
}

func (vs *VisitService) isVisitOpen(user_id int, polygon_id *int) (bool, error) {
	row := db.Db.QueryRow(`SELECT time_exit
								FROM visits vs
								where vs.user_id = &1 and vs.polygon_id = $2 
								order by time_entry desc
								limit 1`,
		user_id, polygon_id)

	//SELECT time_exit is not null
	//FROM visits vs
	//where vs.user_id = 1 and vs.polygon_id = 1
	//order by time_entry desc
	//limit 1

	var timeExit time.Time
	err := row.Scan(&timeExit)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, errors.New("err in check is visit open" + err.Error())
	}

	return true, nil
}
