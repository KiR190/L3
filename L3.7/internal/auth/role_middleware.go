package auth

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

func RequireAnyRole(roles ...string) ginext.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *ginext.Context) {
		role, ok := GetUserRole(c.Request.Context())
		if !ok || role == "" {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if _, ok := allowed[role]; !ok {
			c.JSON(http.StatusForbidden, ginext.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
