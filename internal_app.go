package bonree

import (
	"net/http"
	"bonree/sdk"
)

type app struct {
	appHandle sdk.AppHandle
}

func (app *app) StartBusinessTransaction(name string, w http.ResponseWriter, r *http.Request) BusinessTransaction {
	return newBtn(app, name, w, r)
}

func (app *app) Release() {
	sdk.AppRelease(app.appHandle)
}

// NewApp create application
func NewApp(handle sdk.AppHandle) Application {
	return &app{
		appHandle: handle,
	}
}