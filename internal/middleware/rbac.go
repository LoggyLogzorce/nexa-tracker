package middleware

import (
	"context"
	"errors"
	"net/http"
	"nexa-task-tracker/internal/ctxkeys"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/core/project"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement RBAC
		// 1. Get user role from context
		_, exists := c.Get(ctxkeys.UserRoleKey)
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

func RequireProjectAccess(projectRepo project.Repository, participantRepo participant.Repository, minRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get project_id from URL params
		projectIDStr := c.Param("id")
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
			c.Abort()
			return
		}

		// 2. Get user_id from context
		userIDVal, exists := c.Get(ctxkeys.UserIDKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			c.Abort()
			return
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// 3. Check if project exists and get owner
		proj, err := projectRepo.GetByID(ctx, projectID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		var userRole string

		// 4. Check if user is owner
		if proj.OwnerID == userID {
			userRole = "owner"
		} else {
			// 5. Check if user is participant
			participant, err := participantRepo.GetByProjectAndUser(projectID, userID.String())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
					c.Abort()
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				c.Abort()
				return
			}
			userRole = participant.Role
		}

		// 6. Check if user has required role
		if !hasRequiredRole(userRole, minRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		// 7. Set project role in context
		c.Set(ctxkeys.ProjectRoleKey, userRole)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// hasRequiredRole проверяет, имеет ли пользователь достаточные права
// Иерархия: owner > member > read_only
func hasRequiredRole(userRole, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"owner":     3,
		"member":    2,
		"read_only": 1,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}
