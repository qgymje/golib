package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Check() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodHead && ctx.Request.URL.String() == "/check" {
			ctx.String(200, "im fine")
			ctx.Abort()
		}
	}
}
