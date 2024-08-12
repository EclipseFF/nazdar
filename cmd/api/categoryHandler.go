package main

import (
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
