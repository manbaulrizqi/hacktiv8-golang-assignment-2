package routers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"hacktiv8-golang-assignment-2/api"
	"hacktiv8-golang-assignment-2/datas"
	"net/http"
	"net/mail"
	"strings"
)

func UseApiGroup(e *echo.Group) {
	e.POST("/register", func(c echo.Context) error {

		username := c.FormValue("username")
		password := c.FormValue("password")
		first_name := c.FormValue("first_name")
		last_name := c.FormValue("last_name")

		if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" || strings.TrimSpace(first_name) == "" || strings.TrimSpace(last_name) == "" {
			return c.JSON(http.StatusBadRequest, api.ApiResponse{
				"status":  "fail",
				"message": "Please fill all form input",
			})
		} else {
			_, err := mail.ParseAddress(username)
			if err != nil {
				return c.JSON(http.StatusBadRequest, api.ApiResponse{
					"status":  "fail",
					"message": "Please fill username with email",
				})
			}
		}

		// Get By username
		var allusers []datas.Users
		err := datas.Db.Where("username = ?", username).Limit(1, 0).Find(&allusers)

		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ApiResponse{
				"status":  "fail",
				"message": "Username not found",
			})
		} else if len(allusers) > 0 {
			return c.JSON(http.StatusBadRequest, api.ApiResponse{
				"status":  "fail",
				"message": "Username " + username + " already exists",
			})
		}

		row_user := new(datas.Users)
		row_user.Username = username
		row_user.Password = password
		row_user.FirstName = first_name
		row_user.LastName = last_name

		affected, err := datas.Db.Insert(row_user)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.ApiResponse{
				"status":  "fail",
				"message": "Error " + err.Error(),
			})
		} else {
			fmt.Println("affected:", affected)

			return c.JSON(http.StatusOK, api.ApiResponse{
				"status":  "OK",
				"message": "Register success",
			})
		}
	})
}
