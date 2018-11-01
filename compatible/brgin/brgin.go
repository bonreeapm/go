package brgin

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"bonree"
)

type headerResponseWriter struct{ w gin.ResponseWriter }

func (w *headerResponseWriter) Header() http.Header       { return w.w.Header() }
func (w *headerResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (w *headerResponseWriter) WriteHeader(int)           {}

type replacementResponseWriter struct {
	gin.ResponseWriter
	btn     bonree.BusinessTransaction
	code    int
	written bool
}

var (
	ctxKey = "bonreeBusinessTransaction"
)

func GetCurrentTransaction(c *gin.Context) bonree.BusinessTransaction {
	if v, exists := c.Get(ctxKey); exists {
		if btn, ok := v.(bonree.BusinessTransaction); ok {
			return btn
		}
	}
	return nil
}

func Middleware(app bonree.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.HandlerName()
		// w := &headerResponseWriter{w: c.Writer}
		btn := app.StartBusinessTransaction(name, c.Writer, c.Request)
		defer btn.End()

		// c.Writer = &replacementResponseWriter{
		// 	ResponseWriter: c.Writer,
		// 	btn:            btn,
		// 	code:           http.StatusOK,
		// }
		c.Set(ctxKey, btn)
		c.Next()
	}
}