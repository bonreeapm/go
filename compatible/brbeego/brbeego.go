package brbeego

import (
	"bonree"
	
	"net/http"
	"io"
	"log"
)

type bodyWrapper struct {
	Body   io.ReadCloser
	btn bonree.BusinessTransaction
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

func wrapHandle(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if bw, ok := r.Body.(*bodyWrapper); ok {
			handler.ServeHTTP(bw.btn, r)
		}else {
			handler.ServeHTTP(w, r)
		}

		unWrapRequest(r)
	})
}

func WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, http.Handler) {
	p, h := wrapHandle(pattern, http.HandlerFunc(handler))
	return p, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) })
}

func WrapResponseWriter(app bonree.Application, pattern string, w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	log.Println("StartBusinessTransaction")
	btn := app.StartBusinessTransaction(pattern, w, r)

	wrapRequest(r, btn)

	return btn
}

func GetCurrentTransaction(w http.ResponseWriter) bonree.BusinessTransaction {
	log.Println("GetCurrentTransaction")
	if btn, ok := w.(bonree.BusinessTransaction); ok {
		return btn
	}

	return nil
}