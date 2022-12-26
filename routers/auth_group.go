package routers

import (
	"github.com/labstack/echo/v4"
	"hacktiv8-golang-assignment-2/api"
	"hacktiv8-golang-assignment-2/datas"
	"net/http"
	"os"
	"strings"
)

func UseAuthGroup(e *echo.Group) {
	auth := e.Group("/auth")
	auth.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
			return c.JSON(http.StatusBadRequest, api.ApiResponse{
				"status":  "fail",
				"message": "Please fill username and password",
			})
		}

		// Get By username
		var allusers []datas.Users
		err := datas.Db.Where("username = ?", username).Limit(1, 0).Find(&allusers)

		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ApiResponse{
				"status":  "fail",
				"message": "Username not found",
			})
		} else if len(allusers) == 0 {
			return c.JSON(http.StatusBadRequest, api.ApiResponse{
				"status":  "fail",
				"message": "Username not found",
			})
		} else if allusers[0].Password != password {
			return c.JSON(http.StatusBadRequest, api.ApiResponse{
				"status":  "fail",
				"message": "Password not match",
			})
		}

		// get session by os.Getenv("SESSION_KEY")
		session, _ := api.Store.Get(c.Request(), os.Getenv("SESSION_KEY"))

		session.Values["username"] = allusers[0].Username
		session.Values["first_name"] = allusers[0].FirstName
		session.Values["last_name"] = allusers[0].LastName

		// store session
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, api.ApiResponse{
			"status":  "OK",
			"message": "Login success",
		})
	})
	auth.POST("/logout", func(c echo.Context) error {
		// process get session
		session, _ := api.Store.Get(c.Request(), os.Getenv("SESSION_KEY"))

		// process to expired session
		session.Options.MaxAge = -1
		_ = session.Save(c.Request(), c.Response())

		return c.JSON(http.StatusOK, api.ApiResponse{
			"status":  "OK",
			"message": "Logout success",
		})
	})
	auth.GET("/", func(c echo.Context) error {
		// process get session
		session, _ := api.Store.Get(c.Request(), os.Getenv("SESSION_KEY"))

		// validate session is no data
		if len(session.Values) == 0 {
			return c.JSON(http.StatusOK, api.ApiResponse{
				"status": "fail",
				"error":  "no sessions",
			})
		}

		return c.JSON(http.StatusOK, api.ApiResponse{
			"status": "OK",
			"data": api.ApiResponse{
				"username":   session.Values["username"],
				"first_name": session.Values["first_name"],
				"last_name":  session.Values["last_name"],
			},
		})
	})

}
