package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"nexa-task-tracker/internal/core/attachment"
	"nexa-task-tracker/internal/core/auth"
	"nexa-task-tracker/internal/core/comment"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/task"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/middleware"
)

type Handlers struct {
	AuthHdl        *auth.Handler
	UserHdl        *user.Handler
	ProjectHdl     *project.Handler
	StatusHdl      *status.Handler
	PriorityHdl    *priority.Handler
	ParticipantHdl *participant.Handler
	TaskHdl        *task.Handler
	CommentHdl     *comment.Handler
	AttachmentHdl  *attachment.Handler
}

type Router struct {
	handlers        Handlers
	engine          *gin.Engine
	jwtSecret       string
	corsOrigins     []string
	uploadPath      string
	projectRepo     project.Repository
	participantRepo participant.Repository
	taskRepo        task.Repository
}

func NewRouter(h Handlers, jwtSecret string, corsOrigins []string, uploadPath string, projectRepo project.Repository, participantRepo participant.Repository, taskRepo task.Repository) *Router {
	return &Router{
		handlers:        h,
		engine:          gin.Default(),
		jwtSecret:       jwtSecret,
		corsOrigins:     corsOrigins,
		uploadPath:      uploadPath,
		projectRepo:     projectRepo,
		participantRepo: participantRepo,
		taskRepo:        taskRepo,
	}
}

func (r *Router) Setup() *gin.Engine {
	// CORS
	r.engine.Use(cors.New(cors.Config{
		AllowOrigins:     r.corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// Create rate limiter (5 requests per minute)
	authRateLimiter := middleware.CreateRateLimiter(5)

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve uploaded files as static (avatars, etc.)
	r.engine.Static("./uploads", r.uploadPath)

	// API v1
	v1 := r.engine.Group("/api/v1")
	v1.Use(middleware.BodySizeLimit(5 << 20))
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
				users.PUT("/me/avatar", r.handlers.UserHdl.UploadAvatar)
				users.DELETE("/me", r.handlers.UserHdl.DeleteMe)
				users.GET("/search", r.handlers.UserHdl.SearchUsers)
			}

			protected.GET("/tasks/me", r.handlers.TaskHdl.GetByUserID)

			// Project routes
			projects := protected.Group("/projects")
			{
				projects.GET("", r.handlers.ProjectHdl.List)
				projects.GET("/owned", r.handlers.ProjectHdl.ListOwned)
				projects.POST("", r.handlers.ProjectHdl.Create)

				// Routes requiring project access - read operations (read_only+)
				projectAccess := projects.Group("/:id")
				projectAccess.Use(middleware.RequireProjectAccess(r.projectRepo, r.participantRepo, "read_only"))
				{
					projectAccess.GET("", r.handlers.ProjectHdl.GetByID)
					projectAccess.GET("/participants", r.handlers.ParticipantHdl.GetByProjectID)
					projectAccess.GET("/statuses", r.handlers.StatusHdl.GetByProjectID)
					projectAccess.GET("/priorities", r.handlers.PriorityHdl.GetByProjectID)
					projectAccess.GET("/attachments", r.handlers.AttachmentHdl.GetByProjectID)

					tasks := projectAccess.Group("/tasks")
					tasks.Use(middleware.CheckTaskProject(r.taskRepo))
					{
						tasks.GET("", r.handlers.TaskHdl.GetByProjectID)
						tasksAccess := tasks.Group("/:task_id")
						{
							tasksAccess.GET("", r.handlers.TaskHdl.GetByID)
							tasksAccess.GET("/history", r.handlers.TaskHdl.GetHistoryByTaskID)

							tasksAccess.GET("/attachments", r.handlers.AttachmentHdl.GetByTaskID)
							tasksAccess.GET("/attachments/:attachment_id", r.handlers.AttachmentHdl.Download)
						}
					}
				}

				// Write operations requiring member role
				projectMember := projects.Group("/:id")
				projectMember.Use(middleware.RequireProjectAccess(r.projectRepo, r.participantRepo, "member"))
				{
					projectMember.POST("/statuses", r.handlers.StatusHdl.Create)
					projectMember.PUT("/statuses/:status_id", r.handlers.StatusHdl.Update)
					projectMember.POST("/priorities", r.handlers.PriorityHdl.Create)
					projectMember.PUT("/priorities/:priority_id", r.handlers.PriorityHdl.Update)

					// FIX: создаём свою группу tasks, а не наследуем от projectAccess
					tasks := projectMember.Group("/tasks")
					tasks.Use(middleware.CheckTaskProject(r.taskRepo))
					{
						tasks.POST("", r.handlers.TaskHdl.Create)

						tasksMember := tasks.Group("/:task_id")
						{
							tasksMember.PUT("", r.handlers.TaskHdl.Update)
							tasksMember.GET("/comments", r.handlers.CommentHdl.GetByTaskID)
							tasksMember.POST("/comments", r.handlers.CommentHdl.Create)
							tasksMember.PUT("/comments/:comment_id", r.handlers.CommentHdl.Update)
							tasksMember.DELETE("/comments/:comment_id", r.handlers.CommentHdl.Delete)
							tasksMember.POST("/attachments", r.handlers.AttachmentHdl.Upload)
							tasksMember.DELETE("/attachments/:attachment_id", r.handlers.AttachmentHdl.Delete)
						}
					}
				}

				// Owner-only operations
				projectOwner := projects.Group("/:id")
				projectOwner.Use(middleware.RequireProjectAccess(r.projectRepo, r.participantRepo, "owner"))
				{
					projectOwner.PUT("", r.handlers.ProjectHdl.Update)
					projectOwner.DELETE("", r.handlers.ProjectHdl.Delete)

					// FIX: projectOwner.POST, а не projectMember
					projectOwner.POST("/participants", r.handlers.ParticipantHdl.AddParticipant)
					projectOwner.PUT("/participants/:user_id", r.handlers.ParticipantHdl.UpdateRole)
					projectOwner.DELETE("/participants/:user_id", r.handlers.ParticipantHdl.RemoveParticipant)

					projectOwner.DELETE("/statuses/:status_id", r.handlers.StatusHdl.Delete)
					projectOwner.DELETE("/priorities/:priority_id", r.handlers.PriorityHdl.Delete)

					// FIX: явное указание метода и мидлвара
					projectOwner.DELETE("/tasks/:task_id", middleware.CheckTaskProject(r.taskRepo), r.handlers.TaskHdl.Delete)
				}
			}
		}
	}

	return r.engine
}
