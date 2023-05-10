package auth

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	signingMethodErrorMsg = "invalid signing method"
	tokenClaimsErrorMsg   = "token claims are not of type *tokenClaims"
)

// extractUserIdFromToken parses the provided access token and returns the user ID if successful
func extractUserIdFromToken(accessToken string) (*uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(signingMethodErrorMsg)
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New(tokenClaimsErrorMsg)
	}

	return &claims.UserId, nil
}

// TODO: to be added associations and checking jwt
// GetUserId tries to get the user id from context and return it if successful
func GetUserId(c *gin.Context) (*uint, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return nil, nil
	}

	idInt, ok := id.(uint)
	if !ok {
		return nil, fmt.Errorf("user ID is of invalid type: %v", id)
	}

	return &idInt, nil
}
