package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *App) CheckAdminSession(c echo.Context) error {
	req := struct {
		Token *string `json:"token"`
	}{}
	if err := c.Bind(&req); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	fmt.Println(req.Token)
	session, err := app.repos.Session.GetAdminSessionByToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, session)
}
