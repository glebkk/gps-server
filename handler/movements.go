package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gps_api/db"
	"gps_api/middleware"
	"gps_api/model"
	"net/http"
	"strconv"
	"time"
)

type MovementsHandler struct {
}

func (mh *MovementsHandler) GetAllById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "invalid id")
		return
	}

	var timeEntry = TimeEntry{}
	err = ctx.ShouldBind(&timeEntry)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "invalid time")
		return
	}
	var movements []model.Movement
	var rows *sql.Rows

	if timeEntry.TimeStart.IsZero() {
		rows, err = db.Db.Query("select * from user_movements where user_id=$1", id)
	} else if timeEntry.TimeEnd.IsZero() {
		rows, err = db.Db.Query("select * from user_movements where user_id=$1 and created_at between $2 and now()::timestamp", id, timeEntry.TimeStart)
	} else {
		rows, err = db.Db.Query("select * from user_movements where user_id=$1 and created_at between $2 and $3", id, timeEntry.TimeStart, timeEntry.TimeEnd)
	}
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "err")
		return
	}
	defer rows.Close()
	var movement = model.Movement{}

	for rows.Next() {
		err := rows.Scan(&movement.Id, &movement.UserId, &movement.Latitude, &movement.Longitude, &movement.CreatedAt)
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithStatusJSON(400, "err")
			return
		}
		movements = append(movements, movement)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "err")
		return
	}
	ctx.JSON(http.StatusOK, movements)

}

func (mh *MovementsHandler) GetAll(ctx *gin.Context) {
	var movements []model.Movement
	rows, err := db.Db.Query("select * from user_movements")
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "err")
		return
	}
	defer rows.Close()
	var movement = model.Movement{}

	for rows.Next() {
		err := rows.Scan(&movement.Id, &movement.UserId, &movement.Latitude, &movement.Longitude, &movement.CreatedAt)
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithStatusJSON(400, "err")
			return
		}
		movements = append(movements, movement)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "err")
		return
	}
	ctx.JSON(http.StatusOK, movements)

}

func (mh *MovementsHandler) AddMovement(ctx *gin.Context) {
	id := ctx.GetInt(middleware.UserIdContextVar)

	body := model.MovementCreate{}
	raw, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(404, "invalid data")
		return
	}
	fmt.Println(raw)
	err = json.Unmarshal(raw, &body)
	fmt.Println(body.Longitude, body.Latitude)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	//use Exec whenever we want to insert update or delete
	//Doing Exec(query) will not use a prepared statement, so lesser TCP calls to the SQL server
	_, err = db.Db.Exec(`insert into user_movements ("user_id","latitude", "longitude") values ($1,$2,$3)`, id, body.Latitude, body.Longitude)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create the new movement")
	} else {
		ctx.JSON(http.StatusOK, "Movement created")
	}
}

type TimeEntry struct {
	TimeStart time.Time `form:"timeStart" time_format:"2006-01-02T15:04:05Z"`
	TimeEnd   time.Time `form:"timeEnd" time_format:"2006-01-02T15:04:05Z"`
}
