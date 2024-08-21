package main

import (
	"context"
	"errors"
	"flowers/internal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	port string
	dsn  string
}
type Repos struct {
	Item     *internal.ItemRepo
	User     *internal.UserRepo
	Session  *internal.SessionRepo
	Category *internal.CategoryRepo
	Cart     *internal.CartRepo
	Admin    *internal.AdminRepo
}

type App struct {
	echo   *echo.Echo
	config Config
	pool   *pgxpool.Pool
	repos  *Repos
}

func main() {

	pgDsn := "postgres://postgres:asd123@localhost:5432/nazdar?sslmode=disable"
	app := App{config: Config{port: ":4000", dsn: pgDsn}}
	pool, err := ConnectDB(app.config.dsn)
	app.pool = pool
	if err != nil {
		log.Fatal(err)
	}
	app.echo = echo.New()

	app.repos = &Repos{
		Item:     &internal.ItemRepo{Pool: pool},
		User:     &internal.UserRepo{Pool: pool},
		Session:  &internal.SessionRepo{Pool: pool},
		Category: &internal.CategoryRepo{Pool: pool},
		Cart:     &internal.CartRepo{Pool: pool},
		Admin:    &internal.AdminRepo{Pool: pool},
	}
	app.UseMiddleware()
	app.AddRoutes()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := app.echo.Start(app.config.port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.echo.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.echo.Shutdown(ctx); err != nil {
		app.echo.Logger.Fatal(err)
	}
}
