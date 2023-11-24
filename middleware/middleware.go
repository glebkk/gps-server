package middleware

import (
	"fmt"
	"gps_api/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

const TokenContextVar = "token"
const UserIdContextVar = "user_id"

// TokenValidator returns a gin.HanlerFunc which validates the JWT signature
// and sets a context variable containing the claims.
func CheckToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorization := context.GetHeader("Authorization")
		fmt.Println(authorization)
		row := db.Db.QueryRow(`SELECT "uuid", "id" FROM "users" WHERE "uuid"=$1`, authorization)
		var uuid string
		var id int
		err := row.Scan(&uuid, &id)
		if err != nil {
			fmt.Println(err)
			context.AbortWithStatusJSON(http.StatusUnauthorized, "Nope")
			return
		}
		context.Set(TokenContextVar, authorization)
		context.Set(UserIdContextVar, id)
		context.Next()
	}
}
