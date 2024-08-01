package main

func (app *App) AddRoutes() {
	api := app.echo.Group("/api")
	version := api.Group("/v1")

	version.Static("/images", "./static/")

	items := version.Group("/item")
	items.GET("/:id", app.ReadItemById)
	items.GET("", app.ReadItemsPagination)
	items.POST("", app.CreateItem)
}
