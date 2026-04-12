package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	_ "os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"nexa-task-tracker/internal/api"
	"nexa-task-tracker/internal/config"
	"nexa-task-tracker/internal/core/attachment"
	"nexa-task-tracker/internal/core/auth"
	"nexa-task-tracker/internal/core/comment"
	"nexa-task-tracker/internal/core/history"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/task"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/db"
	"nexa-task-tracker/internal/pkg/events"
)

func main() {
	gin.SetMode(gin.DebugMode)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	database, err := db.Connect(db.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.Migrate(database,
		&user.User{},
		&auth.RefreshToken{},
		&project.Project{},
		&participant.ProjectParticipant{},
		&status.Status{},
		&priority.Priority{},
		&task.Task{},
		&comment.Comment{},
		&history.UpdateHistory{},
		&attachment.Attachment{},
	); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Create EventBus
	eventBus := events.NewEventBus()

	// Initialize repositories
	userRepo := user.NewRepository(database)
	authRepo := auth.NewRepository(database)
	projectRepo := project.NewRepository(database)
	statusRepo := status.NewRepository(database)
	priorityRepo := priority.NewRepository(database)
	participantRepo := participant.NewRepository(database)

	// Initialize services
	userService := user.NewService(userRepo, eventBus)
	authService := auth.NewService(authRepo, userRepo, cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	projectService := project.NewService(projectRepo, eventBus, participantRepo)
	statusService := status.NewService(statusRepo)
	priorityService := priority.NewService(priorityRepo)

	// Initialize handlers
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService, cfg.Cookie.Domain, cfg.Cookie.SameSite, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	projectHandler := project.NewHandler(projectService)
	statusHandler := status.NewHandler(statusService, eventBus)
	priorityHandler := priority.NewHandler(priorityService, eventBus)

	h := api.Handlers{
		AuthHdl:     authHandler,
		UserHdl:     userHandler,
		ProjectHdl:  projectHandler,
		StatusHdl:   statusHandler,
		PriorityHdl: priorityHandler,
	}

	// Setup router
	router := api.NewRouter(h, cfg.JWT.Secret, projectRepo, participantRepo)
	engine := router.Setup()

	// Setup modules
	//if cfg.Modules.Notify {
	//	notify.Init(database, eventBus, engine, cfg.JWT.Secret)
	//}

	// Start server
	go func() {
		if err := engine.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
