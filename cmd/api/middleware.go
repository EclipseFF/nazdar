package main

import "github.com/labstack/echo/v4/middleware"

func (app *App) UseMiddleware() {
	DefaultLoggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}",` +
			`"method":"${method}","uri":"${uri}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}
	app.echo.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
	app.echo.Use(middleware.Recover())
	app.echo.Use(middleware.CORS())
	//app.echo.Use(middleware.CSRF())
	app.echo.Use(middleware.Secure())

	//app.echo.Use(middleware.Static("/static"))
}
