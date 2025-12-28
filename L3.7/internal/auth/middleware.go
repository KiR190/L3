package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/wb-go/wbf/ginext"
)

// ContextKey type for context keys
type ContextKey string

const (
	UserIDKey    ContextKey = "user_id"
	UserEmailKey ContextKey = "user_email"
	UserRoleKey  ContextKey = "user_role"
)

// RequireAuth is middleware that validates JWT tokens
func RequireAuth(authService *AuthService) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := authService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": err.Error()})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
