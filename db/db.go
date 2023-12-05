package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // don't forget to add it. It doesn't be added automatically
	"os"
	"strconv"
)

var Db *sql.DB //created outside to make it global.

// make sure your function start with uppercase to call outside of the directory.
func ConnectDatabase() {

	err := godotenv.Load() //by default, it is .env so we don't have to write
	if err != nil {
		panic("Error is occurred  on .env file please check")
	}
	//we read our .env file
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT")) // don't forget to convert int since port is int type.
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	//var pass string
	//if value, ok := os.LookupEnv("PASSWORD"); ok {
	//	pass = value
	//} else {
	//	pass = ""
	//}

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		fmt.Errorf("There is an error while connecting to the database")
		panic(err)
	}
	if err = db.Ping(); err != nil {
		fmt.Errorf("There is an error while connecting to the database")
		panic(err)
	}

	Db = db
	fmt.Println("Successfully connected to database!")
}
