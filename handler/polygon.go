package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gps_api/model"
	"gps_api/services"
)

type PolygonHandler struct {
	polygonService *services.PolygonService
}

func NewPolygonHandler(polygonService *services.PolygonService) *PolygonHandler {
	return &PolygonHandler{
		polygonService: polygonService,
	}
}

func (ph *PolygonHandler) CreatePolygon(ctx *gin.Context) {
	body := model.PolygonCreate{}
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
	fmt.Println(body.Geometry)
	err = ph.polygonService.CreatePolygon(body)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, err)
		return
	}
	ctx.JSON(200, body)
}
func (ph *PolygonHandler) GetAll(ctx *gin.Context) {
	polygons, err := ph.polygonService.GetAll()
	fmt.Println(polygons)
	if err != nil {
		ctx.AbortWithStatusJSON(500, "error get all")
		return
	}
	ctx.JSON(200, polygons)
}
