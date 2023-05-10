package errorutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func Error(c *gin.Context, code int, message, details string) {
	c.AbortWithStatusJSON(code, ErrorResponse{Status: code, Message: message, Details: details})
}

func InternalServer(c *gin.Context, message, details string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: http.StatusInternalServerError, Message: message, Details: details})
}

func BadRequest(c *gin.Context, message, details string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: http.StatusBadRequest, Message: message, Details: details})
}

func NotFound(c *gin.Context, message, details string) {
	c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Status: http.StatusNotFound, Message: message, Details: details})
}

func Unauthorized(c *gin.Context, message, details string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Status: http.StatusUnauthorized, Message: message, Details: details})
}
