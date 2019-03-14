package main

import (
	"net/http"
	"log"
	"github.com/bonreeapm/go"
	"github.com/bonreeapm/go/compatible/brecho"
	"github.com/labstack/echo"
)

func welcome(c echo.Context) error {
	btn := bonree.GetCurrentTransaction(c.Response().Writer)

	if btn == nil {
		return c.String(http.StatusOK, "Get Transaction fail")
	}

	log.Println("Hello, World!\n")
	return c.String(http.StatusOK, "Hello, World!\n")
}

func main() {
	app, err := bonree.NewApplication("brecho")
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer app.Release()

	e := echo.New()

	e.GET(brecho.WrapHandleFunc(app, "/", welcome))

	e.Logger.Fatal(e.Start(":9099"))
}