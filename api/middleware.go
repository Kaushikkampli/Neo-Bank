package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kaushikkampli/neobank/token"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "authorization"
	AuthorizationType   = "Bearer"
	AuthPayload         = "payload"
)

func authMiddleware(token token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			err := errors.New("authorization header is empty")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 {
			err := errors.New("invalid authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := authFields[0]
		if authType != AuthorizationType {
			err := errors.New("invalid authorization type")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		headerToken := authFields[1]
		payload, err := token.ValidateToken(headerToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(AuthPayload, payload)
		c.Next()
	}
}
