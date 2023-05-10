package auth

import (
	"strings"

	"github.com/emPeeGee/raffinance/pkg/errorutil"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func HandleUserIdentity(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			logger.Info("Unauthenticated request")
			errorutil.Unauthorized(c, "empty auth header", "")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			errorutil.Unauthorized(c, "invalid auth header", "")
			return
		}

		if len(headerParts[1]) == 0 {
			errorutil.Unauthorized(c, "the token is empty", "")
			return
		}

		userId, err := extractUserIdFromToken(headerParts[1])
		if err != nil {
			errorutil.Unauthorized(c, "the token is invalid", err.Error())
			return
		}

		c.Set(userCtx, *userId)
	}
}
