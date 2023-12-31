package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gps_api/db"
	"gps_api/middleware"
	"gps_api/model"
	"gps_api/services"
	"gps_api/valhalla"
	"gps_api/ws"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MovementsHandler struct {
	polygonService *services.PolygonService
	visitService   *services.VisitService
}

func NewMovementHandler(polygonService *services.PolygonService, visitService *services.VisitService) *MovementsHandler {
	return &MovementsHandler{
		polygonService: polygonService,
		visitService:   visitService,
	}
}

func (mh *MovementsHandler) GetAllById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "invalid id")
		return
	}

	layout := "2006-01-02T15:04"
	var timeStart = time.Time{}
	var timeEnd = time.Time{}
	timeStart, err = time.Parse(layout, ctx.DefaultQuery("timeStart", time.Time{}.String()))
	timeEnd, err = time.Parse(layout, ctx.DefaultQuery("timeEnd", time.Time{}.String()))

	fmt.Println(timeStart.Format(time.DateTime), timeStart.IsZero())
	fmt.Println(timeEnd.Format(time.DateTime), timeEnd.IsZero())

	var movements []model.Movement
	var rows *sql.Rows

	if !timeStart.IsZero() && !timeEnd.IsZero() {
		rows, err = db.Db.Query("select * from movements where user_id=$1 and created_at between $2 and $3", id, timeStart.Format(time.DateTime), timeEnd.Format(time.DateTime))
	} else if !timeStart.IsZero() {
		rows, err = db.Db.Query("select * from movements where user_id=$1 and created_at between $2 and now()::timestamp", id, timeStart.Format(time.DateTime))
	} else if !timeEnd.IsZero() {
		rows, err = db.Db.Query("select * from movements where user_id=$1 and created_at <= $2", id, timeEnd.Format(time.DateTime))
	} else {
		rows, err = db.Db.Query("select * from movements where user_id=$1", id)
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(400, "err")
		return
	}
	defer rows.Close()
	var movement = model.Movement{}

	for rows.Next() {
		err := rows.Scan(&movement.Id, &movement.UserId, &movement.Longitude, &movement.Latitude, &movement.CreatedAt)
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
	ctx.AbortWithStatusJSON(http.StatusOK, movements)

}

func (mh *MovementsHandler) GetAll(ctx *gin.Context) {
	var movements []model.Movement
	rows, err := db.Db.Query("select * from movements")
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
	ctx.JSON(200, movements)
	return
}

func (mh *MovementsHandler) CreateMovement(ctx *gin.Context) {
	id := ctx.GetInt(middleware.UserIdContextVar)

	body := model.MovementCreate{}
	raw, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(404, "invalid data")
		return
	}

	err = json.Unmarshal(raw, &body)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}

	valhallaResponse := valhalla.SendLocateRequest(body.Latitude, body.Longitude)
	valhallaPoint := valhalla.GetEdge(valhallaResponse)
	var point *model.MovementCreate

	if valhallaPoint != nil {
		point = valhallaPoint
	} else {
		point = &body
	}

	clients := ws.GetClientsByUserID(id)
	for _, client := range clients {
		var err error
		err = client.Conn.WriteJSON(point)
		if err != nil {
			fmt.Println(err)
			_ = client.Conn.Close()
		}
	}
	ctx.JSON(200, valhallaPoint)

	// Отправьте координаты всем клиентам
	polygon, err := mh.polygonService.GetPolygonByPoint(fmt.Sprintf("%f", point.Longitude), fmt.Sprintf("%f", point.Latitude))
	if err != nil {
		fmt.Println("err", err)
		return
	}
	isOpen, err := mh.visitService.CheckVisitOpen(id)
	if polygon == nil && isOpen {
		if err := mh.visitService.CloseVisit(id); err != nil {
			return
		}
	}
	if polygon != nil && !isOpen {
		if err := mh.visitService.OpenVisit(id, polygon.Id); err != nil {
			return
		}
	}
	return

	_, err = db.Db.Exec(`insert into movements ("user_id", "latitude", "longitude") values ($1,$2,$3)`, id, point.Latitude, point.Longitude)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create the new movement")
		return
	}
	ctx.JSON(http.StatusOK, "Movement created")

}
