package bonree

import (
	"net/http"
)

type api struct {
	application Application
}

// Stop is to call Release() func, to release the app object.
func (api *api) Stop() {
	api.application.Release()
}

// WrapHandleFunc is to wrap the http.handle.
func (api *api) WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := wrapHandle(api.application, pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }
}

func wrapHandle(app Application, pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		btn := app.StartBusinessTransaction(pattern, w, r)
		defer btn.End()

		handler.ServeHTTP(btn, r)
	})
}