package errorutil

import (
	"net/http"

	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/gin-gonic/gin"
)

// Handler creates a middleware that handles panics and errors encountered during HTTP request processing.
func Handler(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			l := logger.With(c.Request.Context())
			errResponse := buildErrorResponse(err)

			l.Errorf("Status %d, Message: %s. Details: %s", errResponse.Status, errResponse.Message, errResponse.Details)

			c.AbortWithStatusJSON(errResponse.Status, errResponse)
		}
	}
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err error) ErrorResponse {
	switch err.(type) {
	case *gin.Error:
		return err.(*gin.Error).Err.(ErrorResponse)
	}

	return ErrorResponse{Status: http.StatusInternalServerError, Message: "Interval Server error", Details: "Error"}
}
