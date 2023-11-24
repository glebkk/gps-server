package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gps_api/db"
	"gps_api/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func (uh *UserHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "invalid id")
		return
	}

	var user = model.User{}
	row := db.Db.QueryRow(`select * from "users" where id=$1;`, id)
	err = row.Scan(&user.Id, &user.Uuid, &user.Name)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(404, "not this id")
		return
	}
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(500, "err row scan")
		return
	}

	ctx.JSON(200, user)

}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	var users []model.User
	rows, err := db.Db.Query(`select * from "users"`)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(500, "err query")
		return
	}
	defer rows.Close()
	var user = model.User{}

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Uuid, &user.Name)
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithStatusJSON(500, "err row scan")
			return
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, users)
}

func (uh *UserHandler) AddUser(ctx *gin.Context) {
	body := model.User{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	//use Exec whenever we want to insert update or delete
	//Doing Exec(query) will not use a prepared statement, so lesser TCP calls to the SQL server
	_, err = db.Db.Exec(`insert into "users" ("uuid","name") values ($1,$2)`, body.Uuid, body.Name)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create the new user.")
	} else {
		ctx.JSON(http.StatusOK, "User is successfully created.")
	}
}
