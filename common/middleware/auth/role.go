package authmiddleware

import (
	"errors"
	"net/http"

	ijwt "github.com/beka-birhanu/common/i_jwt"
	"github.com/gin-gonic/gin"
)

// Constants for context keys
const (
	// ContextUserClaims is the key used to store user claims in the Gin context.
	ContextUserClaims = "userClaims"
)

func Authoriz(jwtService ijwt.Service, hasToBeAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("accessToken")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.Status(http.StatusUnauthorized)
			} else {
				c.Status(http.StatusInternalServerError)
			}
			c.Abort()
			return
		}

		// Decode the token using the JWT service
		claims, err := jwtService.Decode(cookie)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		isAdmin, ok := claims["is_admin"].(bool)
		if !ok || isAdmin != hasToBeAdmin {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}

		// Attach claims to the request context
		c.Set(ContextUserClaims, claims)
		c.Next()
	}
}

