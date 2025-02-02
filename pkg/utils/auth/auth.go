package auth

import (
	"net/http"
	"strings"

	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(attribute string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user") // Assuming user is set in context
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// Cast user data
		currentUser, ok := user.(*types.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user data"})
			c.Abort()
			return
		}

		// Check if user has permission to access the resource
		if !hasPermission(currentUser, attribute) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort()
			return
		}

		// Proceed with the request
		c.Next()
	}
}

func hasPermission(user *types.User, attribute string) bool {
	for _, group := range user.Groups {
		if strings.Contains(attribute, group) {
			return true
		}
	}
	return false
}
