package main

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (app *App) AddAdmin(c echo.Context) error {
	req := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	err := app.repos.Admin.AddAdmin(&req.Name, &req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (app *App) ReadAllAdmins(c echo.Context) error {
	admins, err := app.repos.Admin.GetAllAdmins()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, admins)
}

func (app *App) AdminLogin(c echo.Context) error {
	req := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	admin, err := app.repos.Admin.GetAdminByName(&req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "admin not found"})
	}
	err = bcrypt.CompareHashAndPassword([]byte(*admin.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong password"})
	}
	session, err := app.repos.Session.CreateAdminSession(admin.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, session)
}
