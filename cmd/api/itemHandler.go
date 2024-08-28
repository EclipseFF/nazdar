package main

import (
	"encoding/json"
	"flowers/internal"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"slices"
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
	categoryParam := c.QueryParam("category")
	searchParam := c.QueryParam("search")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "limit is invalid"})
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "offset is invalid"})
	}

	items, err := app.repos.Item.GetItemsPagination(&limit, &offset, &categoryParam, &searchParam)

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
		err := os.Mkdir("./items/"+strconv.Itoa(*item.Id), 0755)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	for _, file := range files {
		go func() {
			src, _ := file.Open()

			defer src.Close()

			dst, _ := os.Create("./items/" + strconv.Itoa(*item.Id) + "/" + file.Filename)

			defer dst.Close()

			io.Copy(dst, src)
		}()
	}
	return c.JSON(http.StatusCreated, item)
}

func (app *App) DeleteItem(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is invalid"})
	}

	err = app.repos.Item.DeleteItem(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (app *App) UpdateItem(c echo.Context) error {
	req := struct {
		Id            *string   `form:"id" json:"id"`
		Name          *string   `form:"name" json:"name"`
		Price         *string   `form:"price" json:"price"`
		Description   *string   `form:"description" json:"description"`
		Images        []*string `form:"images" json:"images"`
		CategoriesIds []*string `form:"categories" json:"categories"`
		OldImages     []string  `form:"oldImages" json:"oldImages"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	var newItem internal.Item
	if req.Id != nil {
		temp, err := strconv.Atoi(*req.Id)
		if err != nil || temp < 1 {

			return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad id"})
		}
		newItem.Id = &temp
	}

	if req.Price != nil {
		temp, err := strconv.Atoi(*req.Price)
		if err != nil || temp < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad price"})
		}
		newItem.Price = &temp
	}

	if req.Description != nil {
		newItem.Description = req.Description
	}

	if req.Name != nil {
		newItem.Name = req.Name
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	if len(form.Value["oldImages"]) > 0 {
		oldImages := form.Value["oldImages"][0]
		err = json.Unmarshal([]byte(oldImages), &req.OldImages)
	}

	files := form.File["newImages"]
	if len(files) > 0 {
		os.Mkdir("./items/"+strconv.Itoa(*newItem.Id), 0755)
	}
	newImages := make([]*string, 0)
	newImages = append(newImages, req.Images...)
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer src.Close()
		dst, err := os.Create("./items/" + strconv.Itoa(*newItem.Id) + "/" + file.Filename)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer dst.Close()
		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		newImages = append(newImages, &file.Filename)
	}

	for _, image := range req.OldImages {
		newImages = append(newImages, &image)
	}

	newItem.Images = newImages

	res, err := app.repos.Item.UpdateItem(&newItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	app.repos.Item.DeleteCategoriesFromItem(res.Id)

	temCats := make([]*string, 0)
	cats := make([]*internal.Category, 0)
	err = json.Unmarshal([]byte(*req.CategoriesIds[0]), &temCats)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	for _, id := range temCats {
		temp, err := strconv.Atoi(*id)
		if err != nil || temp < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad id"})
		}
		cats = append(cats, &internal.Category{Id: &temp})
	}
	newItem.Categories = cats
	err = app.repos.Item.AddCategoriesToItem(res.Id, newItem.Categories)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	go app.DeleteImages(res.Id)
	return c.JSON(http.StatusOK, req)
}

func (app *App) DeleteImages(itemId *int) {
	item, err := app.repos.Item.GetItemById(itemId)
	if err != nil {
		return
	}

	images := make([]string, 0)
	for _, image := range item.Images {
		images = append(images, *image)
	}

	filenames := make([]string, 0)
	dir, err := os.ReadDir("./items/" + strconv.Itoa(*itemId))
	if err != nil {
		return
	}

	for _, file := range dir {
		filenames = append(filenames, file.Name())
	}
	for _, filename := range filenames {
		if !slices.Contains(images, filename) {
			os.Remove("./items/" + strconv.Itoa(*itemId) + "/" + filename)
		}
	}
}
