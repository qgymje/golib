package middleware

import "github.com/gin-gonic/gin"

func Through(mws ...gin.HandlerFunc) func() []gin.HandlerFunc {
	return func() []gin.HandlerFunc {
		return mws
	}
}
