package main

import (
	"flowers/internal"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (app *App) ReadAllCategories(c echo.Context) error {
	cats, err := app.repos.Category.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cats)
}

func (app *App) ReadCategoryById(c echo.Context) error {
	param := c.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cat, err := app.repos.Category.GetCategoryById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, cat)
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

func (app *App) DeleteCategory(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	response, err := app.repos.Category.DeleteCategory(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, response)
}

func (app *App) UpdateCategory(c echo.Context) error {
	req := internal.Category{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	oldCat, err := app.repos.Category.GetCategoryById(req.Id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	req.Id = oldCat.Id
	if *req.Name == "" {
		req.Name = oldCat.Name
	}

	response, err := app.repos.Category.UpdateCategory(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, response)
}
