package main

import (
	"flowers/internal"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (app *App) ReadItemById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is invalid"})
	}

	item, err := app.repos.Item.GetItemById(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, item)
}

func (app *App) ReadItemsPagination(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "limit is invalid"})
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "offset is invalid"})
	}

	items, err := app.repos.Item.GetItemsPagination(&limit, &offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

func (app *App) CreateItem(c echo.Context) error {
	var req internal.Item
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["images"]
	for _, file := range files {
		req.Images = append(req.Images, &file.Filename)
	}

	item, err := app.repos.Item.CreateItem(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(req.Images) > 0 {
		err := os.Mkdir("./static/"+strconv.Itoa(*item.Id), 0755)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	for _, file := range files {
		go func() {
			src, _ := file.Open()

			defer src.Close()

			dst, _ := os.Create("./static/" + strconv.Itoa(*item.Id) + "/" + file.Filename)

			defer dst.Close()

			io.Copy(dst, src)
		}()
	}
	return c.JSON(http.StatusCreated, item)
}
