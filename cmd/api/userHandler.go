package main

import (
	"flowers/internal"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *App) CreateUser(c echo.Context) error {
	var req internal.User
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var user *internal.User

	user, err = app.repos.User.GetUserByPhoneAndName(req.Phone, req.Name)
	if err != nil {
		user, err = app.repos.User.CreateUser(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

	}

	session, err := app.repos.Session.CreateSession(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	response := map[string]any{
		"user":    user,
		"session": session,
	}
	return c.JSON(http.StatusCreated, response)
}
