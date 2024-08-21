package main

func (app *App) AddRoutes() {
	api := app.echo.Group("/api")
	version := api.Group("/v1")

	version.Static("/images", "./items/")

	items := version.Group("/item")
	items.GET("/:id", app.ReadItemById)
	items.GET("", app.ReadItemsPagination)
	items.POST("", app.CreateItem)
	items.DELETE("/delete/:id", app.DeleteItem)

	users := version.Group("/user")
	users.POST("", app.CreateUser)
	users.PUT("", app.UpdateUser)
	users.GET("/:id", app.ReadUserById)
	users.GET("/token/:token", app.ReadUserByToken)
	users.GET("/orders/:token", app.ReadUserOrders)

	sessions := version.Group("/session")
	sessions.POST("/admin/check", app.CheckAdminSession)
	sessions.POST("/admin/login", app.AdminLogin)

	admins := version.Group("/admin")
	admins.GET("", app.ReadAllAdmins)
	admins.POST("", app.AddAdmin)

	categories := version.Group("/category")
	categories.GET("", app.ReadAllCategories)
	categories.POST("", app.CreateCategory)
	categories.DELETE("/:id", app.DeleteCategory)
	categories.GET("/:id", app.ReadCategoryById)
	categories.PUT("", app.UpdateCategory)

	carts := version.Group("/cart")
	carts.GET("/:token", app.ReadCart)
	carts.POST("", app.AddItemToCart)
	carts.DELETE("/remove", app.DeleteItemFromCart)
	carts.POST("/checkout", app.Checkout)
}
