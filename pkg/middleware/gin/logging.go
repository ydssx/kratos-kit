package mgin

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.Query()
		// body:=c.Request.Body

		// Process request
		c.Next()
		code := 200
		if c.Writer != nil {
			code = c.Writer.Status()
		}
		if strings.HasPrefix(path, "/api") {
			log.Infow(
				"operation", path,
				"method", c.Request.Method,
				"args", query,
				"latency", time.Since(startTime).String(),
				"code", code,
			)
		}
	}
}
