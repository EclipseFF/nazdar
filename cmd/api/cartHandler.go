package main

import (
	"flowers/internal"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *App) AddItemToCart(c echo.Context) error {
	req := struct {
		Token *string `json:"token"`
		Item  *struct {
			ItemId      *int    `form:"id" json:"id"`
			Count       *int    `form:"count" json:"quantity"`
			Name        *string `form:"name" json:"name"`
			Price       *int    `form:"price" json:"price"`
			Description *string `form:"description" json:"description"`
			Images      *string `form:"images" json:"image"`
		} `json:"item"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item := internal.CartItem{
		ItemId:      req.Item.ItemId,
		Count:       req.Item.Count,
		Name:        req.Item.Name,
		Price:       req.Item.Price,
		Description: req.Item.Description,
		Images:      []*string{req.Item.Images},
	}

	if *item.Count < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid count"})
	}

	user, err := app.repos.User.GetUserByToken(req.Token)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	err = app.repos.Cart.AddItem(user.Id, &item)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (app *App) ReadCart(c echo.Context) error {
	token := c.Param("token")
	user, err := app.repos.User.GetUserByToken(&token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	items, err := app.repos.Cart.GetCartItemsByUser(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

func (app *App) DeleteItemFromCart(c echo.Context) error {
	req := struct {
		Token  *string `json:"token"`
		ItemId *int    `json:"itemId"`
	}{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	user, err := app.repos.User.GetUserByToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	err = app.repos.Cart.DeleteItem(user.Id, req.ItemId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (app *App) Checkout(c echo.Context) error {
	req := struct {
		Token *string `json:"token"`
		Items []*struct {
			ItemId      *int      `form:"id" json:"id"`
			Count       *int      `form:"count" json:"quantity"`
			Name        *string   `form:"name" json:"name"`
			Price       *int      `form:"price" json:"price"`
			Description *string   `form:"description" json:"description"`
			Images      []*string `form:"images" json:"images"`
		} `json:"items"`
	}{}
	if err := c.Bind(&req); err != nil {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	user, err := app.repos.User.GetUserByToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	items := make([]*internal.CartItem, len(req.Items))

	for i, item := range req.Items {
		items[i] = &internal.CartItem{
			ItemId:      item.ItemId,
			Count:       item.Count,
			Name:        item.Name,
			Price:       item.Price,
			Description: item.Description,
			Images:      item.Images,
		}
	}
	err = internal.SendApiReq(user, items)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go app.repos.Cart.ClearCart(user.Id)

	return c.JSON(http.StatusOK, nil)
}
