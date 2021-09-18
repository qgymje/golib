package middleware

import (
	"github.com/gin-gonic/gin"
	"gthub.com/qgymje/golib/provider"
)

func RecoveryLogger(logger provider.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					logger.Error("[Recovery] panic recovered:\n%s", err)
				}
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
