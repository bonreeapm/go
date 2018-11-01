package brecho

import (
	"bonree"
	"net/http"
	"github.com/labstack/echo"
)

func WrapHandleFunc(app bonree.Application, pattern string, handler func(echo.Context) error) (string, func(echo.Context) error) {
	return pattern, (func(c echo.Context) error {
		response := c.Response()
		request := c.Request()

		btn := app.StartBusinessTransaction(pattern, response.Writer, request)

		response.Writer = btn

		return handler(c)
	})
}

func GetCurrentTransaction(w http.ResponseWriter) bonree.BusinessTransaction {
	if btn, ok := w.(bonree.BusinessTransaction); ok {
		return btn
	}

	return nil
}