package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gps_api/db"
	"gps_api/handler"
	"gps_api/middleware"
	"gps_api/services"
	"gps_api/ws"
	"net/http"
)

func main() {
	route := gin.Default()

	myCors := cors.DefaultConfig()
	myCors.AllowAllOrigins = true
	route.Use(cors.New(myCors))

	db.ConnectDatabase()
	route.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	polygonService := services.NewPolygonService()
	polygonHandler := handler.NewPolygonHandler(polygonService)

	visitService := services.NewVisitService()

	userHandler := handler.UserHandler{}
	movementsHandler := handler.NewMovementHandler(polygonService, visitService)

	route.POST("/users", userHandler.AddUser)
	route.GET("/users", userHandler.GetAllUsers)
	route.GET("/users/:id", userHandler.GetById)
	route.POST("/users/enable-track", middleware.CheckToken(), userHandler.EnableTrack)
	route.POST("/users/disable-track", middleware.CheckToken(), userHandler.DisableTrack)

	route.GET("/movements", movementsHandler.GetAll)
	route.POST("/movements", middleware.CheckToken(), movementsHandler.CreateMovement)
	route.GET("/movements/:id", movementsHandler.GetAllById)

	route.POST("/polygons", polygonHandler.CreatePolygon)

	//route.GET("/polygons", func(ctx *gin.Context) {
	//	lat := ctx.DefaultQuery("latitude", "")
	//	long := ctx.DefaultQuery("longitude", "")
	//	polygon, err := polygonService.GetPolygonByPoint(lat, long)
	//	if err != nil {
	//		ctx.AbortWithStatusJSON(400, err.Error())
	//	}
	//	ctx.JSON(200, polygon)
	//})

	route.GET("/polygons", polygonHandler.GetAll)

	route.GET("/ws", ws.HandleWebSocket)

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}
