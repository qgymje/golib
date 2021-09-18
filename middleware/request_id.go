package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qgymje/golib/utility"
)

const format = "20060102150405"

func RequestId(prefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		if prefix == "" {
			prefix = "def"
		}
		requestId := fmt.Sprintf("%s_%s_%d_%d", prefix, time.Now().Format(format), utility.RandInt(10000, 99999), utility.RandInt(10000, 99999))
		query := req.URL.Query()
		query.Set("request_id", requestId)
		req.URL.RawQuery = query.Encode()

		ctx.Next()
	}
}
