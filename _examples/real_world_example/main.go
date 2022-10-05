package main

import (
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := NewHandler(
		NewUserUC(NewStorage()),
		NewTranslater())
	h.Register(e)

	e.Logger.Fatal(e.Start(":8000"))
}
