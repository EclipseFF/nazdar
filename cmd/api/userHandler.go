package main

import (
	"flowers/internal"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
		fmt.Println(err)
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

func (app *App) ReadUserById(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	user, err := app.repos.User.GetUserById(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (app *App) UpdateUser(c echo.Context) error {
	req := internal.User{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	oldUser, err := app.repos.User.GetUserById(req.Id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	user, err := app.repos.User.UpdateUser(&req, oldUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (app *App) ReadUserOrders(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	orders, err := app.repos.User.GetUserOrders(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, orders)
}
