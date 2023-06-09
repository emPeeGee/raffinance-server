package accesslog

import (
	"time"

	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/gin-gonic/gin"
)

// Handler returns a middleware that records an access log message for every HTTP request being processed.
func Handler(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := c.Request.Context()
		ctx = log.WithRequest(ctx, c.Request)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// generate an access log message
		logger.With(ctx, "duration", time.Now().Sub(start).Milliseconds(), "status", c.Writer.Status()).
			Infof("%s %s %s %d", c.Request.Method, c.Request.URL.Path, c.Request.Proto, c.Writer.Status())
	}
}
