package auth

import (
	"context"
	"net/http"

	"event-booker/internal/models"
	"event-booker/internal/repository"

	"github.com/wb-go/wbf/ginext"
)

type AdminMiddleware struct {
	userRepo repository.UserRepository
}

func NewAdminMiddleware(userRepo repository.UserRepository) *AdminMiddleware {
	return &AdminMiddleware{userRepo: userRepo}
}

// RequireAdmin checks if the authenticated user has admin role
func (m *AdminMiddleware) RequireAdmin() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		// Get user ID from context (set by auth middleware)
		userID, ok := GetUserID(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Get user from database
		user, err := m.userRepo.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user is admin
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, ginext.H{"error": "admin access required"})
			c.Abort()
			return
		}

		// Set user role in context for handlers
		ctx := context.WithValue(c.Request.Context(), "userRole", string(user.Role))
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

