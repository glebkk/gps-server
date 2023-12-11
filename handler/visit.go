package handler

import (
	"github.com/gin-gonic/gin"
	"gps_api/services"
	"strconv"
)

type VisitHandler struct {
	visitService *services.VisitService
}

func NewVisitHandler(service *services.VisitService) *VisitHandler {
	return &VisitHandler{visitService: service}
}

func (vh *VisitHandler) GetLastUserVisits(ctx *gin.Context) {
	userID, isSet := ctx.GetQuery("id")
	if !isSet {
		ctx.AbortWithStatusJSON(404, gin.H{"error": "need query key id"})
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(404, gin.H{"error": "id should be number"})
		return
	}

	visits, err := vh.visitService.GetLastUserVisits(id)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "something went wrong"})
	}

	ctx.JSON(200, visits)

}
