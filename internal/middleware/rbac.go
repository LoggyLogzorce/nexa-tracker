package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement RBAC
		// 1. Get user role from context
		_, exists := c.Get(UserRoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// 2. Check if user has required role
		// TODO: Check if userRole is in roles slice

		c.Next()
	}
}

func RequireProjectAccess(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement project-level RBAC
		// 1. Get project_id from URL params
		projectIDStr := c.Param("project_id")
		_, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
			c.Abort()
			return
		}

		// 2. Get user_id from context
		// 3. Check if user has required role in project (via participant service)
		// 4. Call c.Next() or c.Abort()

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}
