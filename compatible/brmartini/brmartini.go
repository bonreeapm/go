package brmartini

import (
	"github.com/bonreeapm/go"

	"io"
	"net/http"

	"github.com/go-martini/martini"
)

type bodyWrapper struct {
	Body io.ReadCloser
	btn  bonree.BusinessTransaction
}

func (b *bodyWrapper) Close() error {
	return b.Body.Close()
}
func (b *bodyWrapper) Read(p []byte) (n int, err error) {
	return b.Body.Read(p)
}

func wrapRequest(r *http.Request, btn bonree.BusinessTransaction) bool {
	if isWraped(r) {
		return false
	}
	body := r.Body
	r.Body = &bodyWrapper{Body: body, btn: btn}
	return true
}

func unWrapRequest(r *http.Request) {
	if p, ok := r.Body.(*bodyWrapper); ok {
		r.Body = p.Body
		p.Body = nil
		p.btn = nil
	}
}

func isWraped(r *http.Request) bool {
	if _, ok := r.Body.(*bodyWrapper); ok {
		return true
	}
	return false
}

func MiddlewareHandle(app bonree.Application) func(c martini.Context, w http.ResponseWriter, r *http.Request) {
	return func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		pattern := r.URL.Path
		btn := app.StartBusinessTransaction(pattern, w, r)

		wrapRequest(r, btn)

		c.Next()

		btn.End()

		unWrapRequest(r)
	}
}

func GetCurrentTransaction(r *http.Request) bonree.BusinessTransaction {
	if bw, ok := r.Body.(*bodyWrapper); ok {
		return bw.btn
	}

	return nil
}
