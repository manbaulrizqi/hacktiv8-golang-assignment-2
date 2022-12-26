package main

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"hacktiv8-golang-assignment-2/api"
	"hacktiv8-golang-assignment-2/datas"
	"hacktiv8-golang-assignment-2/routers"
	"os"
	"time"
	"xorm.io/xorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed load .env : ", err.Error())
		return
	}
	var dbConnection = "mssql"
	var dbHost = os.Getenv("DB_HOST")
	var dbName = os.Getenv("DB_DATABASE")
	var dbUser = os.Getenv("DB_USERNAME")
	var dbPass = os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", dbHost, dbUser, dbPass, dbName)
	datas.Db, err = xorm.NewEngine(dbConnection, dsn)
	if err != nil {
		fmt.Println("dsn: ", dsn)
		fmt.Println("fail connect to database : ", err.Error())
		return
	}
	fmt.Println("database connected")
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err == nil {
		datas.Db.SetTZDatabase(loc)
		datas.Db.SetTZLocation(loc)
	}

	api.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// Init Echo Web Server
	e := echo.New()

	// Handle memory leak
	e.Use(echo.WrapMiddleware(context.ClearHandler))
	// Static Asset
	e.Static("/static", "static/static")
	// Static HTML
	e.File("/", "static/views/index.html")
	e.File("/dashboard", "static/views/dashboard.html")
	e.File("/register", "static/views/register.html")

	api := e.Group("/api")
	routers.UseAuthGroup(api)
	routers.UseApiGroup(api)

	e.Start(":4444")
}
