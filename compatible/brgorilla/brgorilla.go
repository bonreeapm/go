package brgorilla

import (
	"net/http"

	"github.com/gorilla/mux"
	"bonree"
)

type instrumentedHandler struct {
	name string
	app  bonree.Application
	orig http.Handler
}

func (h instrumentedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	txn := h.app.StartBusinessTransaction(h.name, w, r)
	defer txn.End()

	h.orig.ServeHTTP(txn, r)
}

func instrumentRoute(h http.Handler, app bonree.Application, name string) http.Handler {
	if _, ok := h.(instrumentedHandler); ok {
		return h
	}
	return instrumentedHandler{
		name: name,
		orig: h,
		app:  app,
	}
}

func routeName(route *mux.Route) string {
	if nil == route {
		return ""
	}
	if n := route.GetName(); n != "" {
		return n
	}
	if n, _ := route.GetPathTemplate(); n != "" {
		return n
	}
	n, _ := route.GetHostTemplate()
	return n
}

// InstrumentRoutes adds instrumentation to a router.  This must be used after
// the routes have been added to the router.
func InstrumentRoutes(r *mux.Router, app bonree.Application) *mux.Router {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		h := instrumentRoute(route.GetHandler(), app, routeName(route))
		route.Handler(h)
		return nil
	})
	if nil != r.NotFoundHandler {
		r.NotFoundHandler = instrumentRoute(r.NotFoundHandler, app, "NotFoundHandler")
	}
	return r
}
