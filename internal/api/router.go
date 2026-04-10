package api

import (
	"github.com/gin-gonic/gin"
	"nexa-task-tracker/internal/core/auth"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/middleware"
)

type Handlers struct {
	AuthHdl     *auth.Handler
	UserHdl     *user.Handler
	ProjectHdl  *project.Handler
	StatusHdl   *status.Handler
	PriorityHdl *priority.Handler
}

type Router struct {
	handlers  Handlers
	engine    *gin.Engine
	jwtSecret string
}

func NewRouter(h Handlers, jwtSecret string) *Router {
	return &Router{
		handlers:  h,
		engine:    gin.Default(),
		jwtSecret: jwtSecret,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Create rate limiter (5 requests per minute)
	authRateLimiter := middleware.CreateRateLimiter(5)

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes (public) with rate limiting
		authGroup := v1.Group("/auth")
		authGroup.Use(middleware.RateLimitMiddleware(authRateLimiter))
		{
			authGroup.POST("/register", r.handlers.AuthHdl.Register)
			authGroup.POST("/login", r.handlers.AuthHdl.Login)
			authGroup.POST("/refresh", r.handlers.AuthHdl.Refresh)
			authGroup.POST("/logout", r.handlers.AuthHdl.Logout)

			// 2FA routes | не сделано
			twoFA := authGroup.Group("/2fa")
			{
				twoFA.POST("/setup", r.handlers.AuthHdl.Setup2FA)
				twoFA.POST("/verify", r.handlers.AuthHdl.Verify2FA)
				twoFA.POST("/enable", r.handlers.AuthHdl.Enable2FA)
				twoFA.POST("/disable", r.handlers.AuthHdl.Disable2FA)
			}
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.Auth(r.jwtSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", r.handlers.UserHdl.GetMe)
				users.PUT("/me", r.handlers.UserHdl.UpdateMe)
				users.DELETE("/me", r.handlers.UserHdl.DeleteMe)
			}

			// Project routes
			projects := protected.Group("/projects")
			{
				projects.GET("", r.handlers.ProjectHdl.List)
				projects.POST("", r.handlers.ProjectHdl.Create)
				projects.GET("/:id", r.handlers.ProjectHdl.GetByID)
				projects.PUT("/:id", r.handlers.ProjectHdl.Update)
				projects.DELETE("/:id", r.handlers.ProjectHdl.Delete)

				// Project participants
				projects.GET("/:id/participants", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get participants"}) })
				projects.POST("/:id/participants", func(c *gin.Context) { c.JSON(200, gin.H{"message": "add participant"}) })
				projects.PUT("/:id/participants/:user_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "update participant"}) })
				projects.DELETE("/:id/participants/:user_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "remove participant"}) })

				// Project statuses
				projects.GET("/:id/statuses", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get statuses"}) })
				projects.POST("/:id/statuses", func(c *gin.Context) { c.JSON(200, gin.H{"message": "create status"}) })
				projects.PUT("/:id/statuses/:status_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "update status"}) })
				projects.DELETE("/:id/statuses/:status_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "delete status"}) })

				// Project priorities
				projects.GET("/:id/priorities", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get priorities"}) })
				projects.POST("/:id/priorities", func(c *gin.Context) { c.JSON(200, gin.H{"message": "create priority"}) })
				projects.PUT("/:id/priorities/:priority_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "update priority"}) })
				projects.DELETE("/:id/priorities/:priority_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "delete priority"}) })
			}

			// Task routes
			tasks := protected.Group("/tasks")
			{
				tasks.GET("", func(c *gin.Context) { c.JSON(200, gin.H{"message": "list tasks"}) })
				tasks.POST("", func(c *gin.Context) { c.JSON(200, gin.H{"message": "create task"}) })
				tasks.GET("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get task"}) })
				tasks.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "update task"}) })
				tasks.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "delete task"}) })

				// Task history
				tasks.GET("/:id/history", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get task history"}) })

				// Task comments
				tasks.GET("/:id/comments", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get comments"}) })
				tasks.POST("/:id/comments", func(c *gin.Context) { c.JSON(200, gin.H{"message": "create comment"}) })
				tasks.PUT("/:id/comments/:comment_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "update comment"}) })
				tasks.DELETE("/:id/comments/:comment_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "delete comment"}) })

				// Task attachments
				tasks.GET("/:id/attachments", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get attachments"}) })
				tasks.POST("/:id/attachments", func(c *gin.Context) { c.JSON(200, gin.H{"message": "upload attachment"}) })
				tasks.GET("/:id/attachments/:attachment_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "download attachment"}) })
				tasks.DELETE("/:id/attachments/:attachment_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "delete attachment"}) })
			}

		}
	}

	return r.engine
}
