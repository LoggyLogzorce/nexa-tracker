package api

import (
	"github.com/gin-gonic/gin"
	"nexa-task-tracker/internal/core/auth"
	"nexa-task-tracker/internal/core/participant"
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
	handlers        Handlers
	engine          *gin.Engine
	jwtSecret       string
	projectRepo     project.Repository
	participantRepo participant.Repository
}

func NewRouter(h Handlers, jwtSecret string, projectRepo project.Repository, participantRepo participant.Repository) *Router {
	return &Router{
		handlers:        h,
		engine:          gin.Default(),
		jwtSecret:       jwtSecret,
		projectRepo:     projectRepo,
		participantRepo: participantRepo,
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

				// Routes requiring project access - read operations (read_only+)
				projectAccess := projects.Group("/:id")
				projectAccess.Use(middleware.RequireProjectAccess(r.projectRepo, r.participantRepo, "read_only"))
				{
					projectAccess.GET("", r.handlers.ProjectHdl.GetByID)
					projectAccess.GET("/participants", func(c *gin.Context) { c.JSON(200, gin.H{"message": "get participants"}) })
					projectAccess.GET("/statuses", r.handlers.StatusHdl.GetByProjectID)
					projectAccess.GET("/priorities", r.handlers.PriorityHdl.GetByProjectID)
				}

				// Write operations requiring member role
				projectMember := projects.Group("/:id")
				projectMember.Use(middleware.RequireProjectAccess(r.projectRepo, r.participantRepo, "member"))
				{
					projectMember.POST("/participants", func(c *gin.Context) { c.JSON(200, gin.H{"message": "add participant"}) })
					projectMember.POST("/statuses", r.handlers.StatusHdl.Create)
					projectMember.PUT("/statuses/:status_id", r.handlers.StatusHdl.Update)
					projectMember.POST("/priorities", r.handlers.PriorityHdl.Create)
					projectMember.PUT("/priorities/:priority_id", r.handlers.PriorityHdl.Update)
				}

				// Owner-only operations
				projectOwner := projects.Group("/:id")
				projectOwner.Use(middleware.RequireProjectAccess(r.projectRepo, r.participantRepo, "owner"))
				{
					projectOwner.PUT("", r.handlers.ProjectHdl.Update)
					projectOwner.DELETE("", r.handlers.ProjectHdl.Delete)
					projectOwner.PUT("/participants/:user_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "update participant"}) })
					projectOwner.DELETE("/participants/:user_id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "remove participant"}) })
					projectOwner.DELETE("/statuses/:status_id", r.handlers.StatusHdl.Delete)
					projectOwner.DELETE("/priorities/:priority_id", r.handlers.PriorityHdl.Delete)
				}
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
