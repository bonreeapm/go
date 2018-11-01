package bonree

import (
	"net/http"
)

// API is provide the method to use SDK.
type API interface {
	Stop()
	WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request))
}

// NewAPI create the new Api
func NewAPI(app Application) (API, error) {
	api := &api {
		application: app,
	}

	return api, nil
}