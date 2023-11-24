package main

import (
	"gps_api/db"
	"gps_api/handler"
	"gps_api/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	userHandler := handler.UserHandler{}
	movementsHandler := handler.MovementsHandler{}

	route.POST("/users", userHandler.AddUser)
	route.GET("/users", userHandler.GetAllUsers)
	route.GET("/users/:id", userHandler.GetById)

	route.GET("/movements", movementsHandler.GetAll)
	route.GET("/movements/:id", movementsHandler.GetAllById)
	route.POST("/movements", middleware.CheckToken(), movementsHandler.AddMovement)
	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}
