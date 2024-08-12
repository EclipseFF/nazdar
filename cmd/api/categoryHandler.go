package main

import (
	"flowers/internal"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *App) ReadAllCategories(c echo.Context) error {
	cats, err := app.repos.Category.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cats)
}

func (app *App) CreateCategory(c echo.Context) error {
	req := internal.Category{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response, err := app.repos.Category.CreateCategory(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, response)
}
